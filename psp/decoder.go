// Package psp implements a Paint Shop Pro image decoder.
package psp

// https://github.com/GNOME/gimp/blob/2275d4b257e9de36f1ac749e591378e58b348754/plug-ins/common/file-psp.c

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"runtime"
	"strings"
	"time"
)

var (
	fileMagic  = []byte("Paint Shop Pro Image File\n\x1a\x00\x00\x00\x00\x00")
	blockMagic = []byte("~BK\x00")
	chunkMagic = []byte("~FL\x00")
)

type decoder struct {
	r              *bufio.Reader
	versionMinor   uint16
	versionMajor   uint16
	width          int
	height         int
	res            float64
	resMetric      metric
	comp           compression
	colorModel     color.Model
	bitDepth       uint16
	planeCount     uint16
	colorCount     uint32
	grayscale      bool
	totalImageSize uint32
	activeLayer    int32
	layerCount     uint16
	xDataTrnsIndex uint16
	creator        creator
	palette        color.Palette
	tmpBuf         []byte
}

type blockHeader struct {
	id      blockID
	dataLen uint32
	initLen uint32 // Only for major ver <= 3
}

type chunkHeader struct {
	fieldKeyword uint16
	dataLen      uint32
}

type creator struct {
	title            string
	creationDate     time.Time
	modificationDate time.Time
	artist           string
	copyright        string
	description      string
	appID            uint32
	appVersion       uint32
}

type layer struct {
	name                  string
	layerType             layerType
	rect                  image.Rectangle
	savedRect             image.Rectangle
	opacity               byte
	blendingMode          byte
	visible               bool
	transparencyProtected bool
	linkGroupID           byte
	maskRectangle         image.Rectangle
	savedMaskRectangle    image.Rectangle
	maskLinked            bool
	maskDisabled          bool
	invertMaskOnBlend     bool
	blendRangeCount       uint16
	bitmapCount           uint16
	channelCount          uint16
}

// A FormatError reports that the input is not a valid PCX.
type FormatError string

func (e FormatError) Error() string {
	return "psp: invalid format: " + string(e)
}

// An UnsupportedError reports that the variant of the PCX file is not supported.
type UnsupportedError string

func (e UnsupportedError) Error() string {
	return "psp: unsupported variant: " + string(e)
}

func init() {
	image.RegisterFormat("psp", string(fileMagic), Decode, DecodeConfig)
}

// Decode reads a PSP image from r and returns it as an image.Image.
// The type of Image returned depends on the PSP contents.
func Decode(r io.Reader) (img image.Image, err error) {
	defer catchErrors(&err)
	d := newDecoder(r)
	return d.decode(), nil
}

// DecodeConfig returns the color model and dimensions of a PSP image
// without decoding the entire image.
func DecodeConfig(r io.Reader) (config image.Config, err error) {
	defer catchErrors(&err)
	d := newDecoder(r)
	return image.Config{
		ColorModel: d.colorModel,
		Width:      d.width,
		Height:     d.height,
	}, nil
}

func catchErrors(err *error) {
	if r := recover(); r != nil {
		if _, ok := r.(runtime.Error); ok {
			panic(r)
		}
		*err = r.(error)
	}
}

func newDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:      bufio.NewReader(r),
		tmpBuf: make([]byte, 64),
	}
	d.readHeader()
	return d
	// if err == io.EOF {
	// 	err = io.ErrUnexpectedEOF
	// }
}

func (d *decoder) error(err error) {
	panic(err)
}

func (d *decoder) readHeader() {
	d.read(d.tmpBuf[:36])
	if !bytes.Equal(d.tmpBuf[:32], fileMagic) {
		d.error(FormatError("not a PSP file"))
	}
	d.versionMajor = decodeUint16(d.tmpBuf[32:34])
	d.versionMinor = decodeUint16(d.tmpBuf[34:36])
	if d.versionMajor < 3 {
		d.error(UnsupportedError("only major versions >= 3 are supported"))
	}

	var bh blockHeader
	d.readBlockHeader(&bh)
	if bh.id != imageBlock {
		d.error(FormatError("missing general image attributes block"))
	} else if bh.dataLen < 38 || bh.dataLen > 64 {
		d.error(FormatError("invalid length for general image attributes block"))
	}
	d.read(d.tmpBuf[:bh.dataLen])
	buf := d.tmpBuf[:bh.dataLen]
	if d.versionMajor >= 4 {
		buf = buf[4:]
	}
	d.width = int(int32(decodeUint32(buf[0:4])))
	d.height = int(int32(decodeUint32(buf[4:8])))
	d.res = math.Float64frombits(decodeUint64(buf[8:16]))
	d.resMetric = metric(buf[16])
	d.comp = compression(decodeUint16(buf[17:19]))
	d.bitDepth = decodeUint16(buf[19:21])
	d.planeCount = decodeUint16(buf[21:23])
	d.colorCount = decodeUint32(buf[23:27])
	d.grayscale = buf[27] == 1
	d.totalImageSize = decodeUint32(buf[28:32])
	d.activeLayer = int32(decodeUint32(buf[32:36]))
	d.layerCount = decodeUint16(buf[36:38])

	// Validate some values
	switch d.comp {
	case compressionNone, compressionRLE, compressionLZ77:
	default:
		d.error(UnsupportedError(fmt.Sprintf("unsupported compression (%04x)", d.comp)))
	}
	if d.grayscale {
		switch d.bitDepth {
		case 8:
			d.colorModel = color.GrayModel
		case 16:
			d.colorModel = color.Gray16Model
		default:
			d.error(UnsupportedError(fmt.Sprintf("unsupported bit depth %d for grayscale image", d.bitDepth)))
		}
	} else {
		switch d.bitDepth {
		// case 1: // TODO: not sure how to decode this properly
		case 16:
			d.colorModel = color.Gray16Model
		case 8, 24:
			d.colorModel = color.RGBAModel
		case 48, 64:
			d.colorModel = color.RGBA64Model
		default:
			d.error(UnsupportedError(fmt.Sprintf("unsupported bit depth %d", d.bitDepth)))
		}
	}
	fmt.Printf("%+v\n", d)
}

func (d *decoder) decode() image.Image {
	for {
		var bh blockHeader
		d.readBlockHeader(&bh)
		switch bh.id {
		case extendedDataBlock:
			d.decodeExtendedDataBlock(int64(bh.dataLen))
		case creatorBlock:
			d.decodeCreatorBlock(int64(bh.dataLen))
		case colorBlock:
			d.decodeColorBlock(int(bh.dataLen))
		case layerStartBlock:
			img, _ := d.decodeLayers()
			return img
		case compositeImageBankBlock: // TODO
			// length?: uint32
			// number of thumbnails?: uint32
			// sub blocks
			//   block ID 0x11 (len 0x18):
			//     length?: uint32
			//     width?: int32
			//     height?: int32
			//     0x0008: uint16
			//     0x0002: uint16
			//     0x0001: uint16
			//     0x 00 0x01 0x00 0x00 0x01 0x00
			//   block ID 0x09 (len 0x0b36)
			//     0x08 0x00 0x 00 0x00 0x01 0x00 0x01 0x00
			//     sub blocks
			//       block ID 0x02 (len 0x0408)
			//       block ID 0x05 (len 0x0712)
			fallthrough
		default:
			d.skip(int(bh.dataLen))
		}
	}
}

func (d *decoder) decodeColorBlock(ln int) {
	if d.versionMajor >= 4 {
		d.readUint32() // TODO: 0x08 maybe color type/format
	}
	nColors := int(d.readUint32())
	if len(d.tmpBuf) < nColors*4 {
		d.tmpBuf = make([]byte, nColors*4)
	}
	d.read(d.tmpBuf[:nColors*4])
	d.palette = make([]color.Color, nColors)
	for i := 0; i < nColors; i++ {
		d.palette[i] = color.RGBA{
			R: d.tmpBuf[i*4+2],
			G: d.tmpBuf[i*4+1],
			B: d.tmpBuf[i*4],
			A: 255, // the last value isn't actually alpha but rather always 0
		}
	}
}

func (d *decoder) decodeLayers() (image.Image, *layer) {
	var layer layer
	var img image.Image
	var imgRGBA *image.RGBA
	var imgRGBA64 *image.RGBA64
	var imgGray16 *image.Gray16
	var imgPaletted *image.Paletted
	var layerBytes int
	channel := 0
	for {
		var bh blockHeader
		d.readBlockHeader(&bh)
		switch bh.id {
		case layerBlock:
			// headerLen := d.readUint32()
			// println(headerLen)
			if d.versionMajor >= 4 {
				d.readUint32() // length? doesn't really match
				nameLen := d.readUint16()
				layer.name = d.readString(int(nameLen))
			} else {
				layer.name = strings.TrimSpace(d.readString(256))
			}
			layer.layerType = layerType(d.readByte())
			layer.rect = d.readRect()
			layer.savedRect = d.readRect()
			layer.opacity = d.readByte()
			layer.blendingMode = d.readByte()
			layer.visible = d.readByte() != 0
			layer.transparencyProtected = d.readByte() != 0
			layer.linkGroupID = d.readByte()
			layer.maskRectangle = d.readRect()
			layer.savedMaskRectangle = d.readRect()
			layer.maskLinked = d.readByte() != 0
			layer.maskDisabled = d.readByte() != 0
			layer.invertMaskOnBlend = d.readByte() != 0
			layer.blendRangeCount = d.readUint16()
			/*
				TODO:
					blend ranges (4 bytes per range) * 5
						source blend range
						destination blend range
			*/
			d.skip(4 * 2 * 5)
			// TODO: not sure about these versions or what's going on
			if d.versionMajor >= 10 {
				d.skip(5)
				// TODO: not sure how to read or calculate these
				if d.palette != nil {
					layer.channelCount = 1
				} else {
					switch d.bitDepth {
					case 1: // TODO: not sure how to decode this properly
						layer.channelCount = 1
					case 8:
						layer.channelCount = 1
					case 16:
						layer.channelCount = 1
					case 24, 48:
						layer.channelCount = 3
					case 32, 64:
						layer.channelCount = 4
					default:
						d.error(FormatError("unknown channel count"))
					}
				}
			} else if d.versionMajor >= 6 {
				d.skip(9)
				layer.bitmapCount = d.readUint16()
				layer.channelCount = d.readUint16()
			} else if d.versionMajor >= 4 {
				d.skip(4)
				layer.bitmapCount = d.readUint16()
				layer.channelCount = d.readUint16()
			} else {
				layer.bitmapCount = d.readUint16()
				layer.channelCount = d.readUint16()
			}
			fmt.Printf("%+v\n", layer)
			if layer.channelCount == 0 {
				break
			}
			channel = 0
			if d.palette != nil {
				imgPaletted = image.NewPaletted(layer.savedRect, d.palette)
				img = imgPaletted
				layerBytes = layer.savedRect.Dx() * layer.savedRect.Dy()
				if d.bitDepth == 1 {
					layerBytes /= 8
				}
			} else if d.bitDepth == 16 {
				imgGray16 = image.NewGray16(layer.savedRect)
				img = imgGray16
				layerBytes = layer.savedRect.Dx() * layer.savedRect.Dy() * 2
			} else if d.bitDepth == 24 || d.bitDepth == 32 {
				imgRGBA = image.NewRGBA(layer.savedRect)
				img = imgRGBA
				for i := 3; i < len(imgRGBA.Pix); i += 4 {
					imgRGBA.Pix[i] = 255
				}
				layerBytes = layer.savedRect.Dx() * layer.savedRect.Dy()
			} else if d.bitDepth == 48 || d.bitDepth == 64 {
				imgRGBA64 = image.NewRGBA64(layer.savedRect)
				img = imgRGBA64
				for i := 6; i < len(imgRGBA64.Pix); i += 8 {
					imgRGBA64.Pix[i] = 255
					imgRGBA64.Pix[i+1] = 255
				}
				layerBytes = layer.savedRect.Dx() * layer.savedRect.Dy() * 2
			}
		case channelBlock:
			if d.versionMajor >= 4 {
				headerLen := d.readUint32()
				if headerLen != 16 {
					d.error(FormatError("invalid channel block info len"))
				}
			}
			compressedLayerLen := int(d.readUint32())
			uncompressedImageLen := int(d.readUint32())
			bitmapType := bitmapType(d.readUint16())
			channelType := channelType(d.readUint16())
			if bitmapType != dibImage {
				// TODO: ignoring other bitmap types (e.g. mask)
				d.skip(int(bh.dataLen - 4*3 - 2*2))

				channel++
				if channel == int(layer.channelCount) {
					return img, &layer
				}
				break
			}
			fmt.Printf("Channel\n")
			fmt.Printf("\tcompressed layer len = %d\n", compressedLayerLen)
			fmt.Printf("\tuncompressed image len = %d\n", uncompressedImageLen)
			fmt.Printf("\tbitmap type = %s\n", bitmapType)
			fmt.Printf("\tchannel type = %s\n", channelType)

			if cap(d.tmpBuf) < layerBytes {
				d.tmpBuf = make([]byte, layerBytes)
			}
			buf := d.tmpBuf[:layerBytes]

			switch d.comp {
			case compressionLZ77:
				zr, err := zlib.NewReader(io.LimitReader(d.r, int64(compressedLayerLen)))
				if err != nil {
					d.error(err)
				}
				_, err = io.ReadFull(zr, buf)
				zr.Close()
				if err != nil {
					d.error(err)
				}
			case compressionRLE:
				j := 0
				for n := compressedLayerLen; n > 0; n-- {
					if run := int(d.readByte()); run > 128 {
						b := d.readByte()
						n--
						for i := 0; i < run-128; i++ {
							buf[j] = b
							j++
						}
					} else {
						n -= run
						d.read(buf[j : j+run])
						j += run
					}
				}
			case compressionNone:
				d.read(buf)
			}

			if imgRGBA != nil {
				for i := int(channelType) - 1; i < len(imgRGBA.Pix); i += 4 {
					imgRGBA.Pix[i] = buf[i/4]
				}
			} else if imgRGBA64 != nil {
				for i := (int(channelType) - 1) * 2; i < len(imgRGBA64.Pix); i += 8 {
					imgRGBA64.Pix[i] = buf[2*(i/8)+1]
					imgRGBA64.Pix[i+1] = buf[2*(i/8)]
				}
			} else if imgGray16 != nil {
				for i := 0; i < len(buf); i += 2 {
					imgGray16.Pix[i] = buf[i+1]
					imgGray16.Pix[i+1] = buf[i]
				}
			} else {
				if d.bitDepth == 1 {
					for i, b := range buf {
						for j := 0; j < 8; j++ {
							imgPaletted.Pix[i*8+j] = b >> 7
							b <<= 1
						}
					}
				} else {
					imgPaletted.Pix = buf
				}
			}

			channel++
			if channel == int(layer.channelCount) {
				return img, &layer
			}
		case 33:
			// TODO: No idea what this block is (shows up in major version 13). seems to be all zeros
			d.skip(int(bh.dataLen))
			n := int(d.readUint32())
			d.skip(n - 4)
		default:
			d.skip(int(bh.dataLen))
		}
	}
}

func (d *decoder) dump(n int) {
	if cap(d.tmpBuf) < n {
		d.tmpBuf = make([]byte, n)
	}
	d.read(d.tmpBuf[:n])
	fmt.Println(hex.Dump(d.tmpBuf[:n]))
}

func (d *decoder) decodeExtendedDataBlock(totalLen int64) {
	var ch chunkHeader
	for totalLen > 0 {
		d.readChunkHeader(&ch)
		totalLen -= 10 + int64(ch.dataLen)
		switch ch.fieldKeyword {
		case xDataTrnsIndex:
			// TODO
			fallthrough
		default:
			d.skip(int(ch.dataLen))
		}
	}
}

func (d *decoder) decodeCreatorBlock(totalLen int64) {
	var ch chunkHeader
	for totalLen > 0 {
		d.readChunkHeader(&ch)
		totalLen -= 10 + int64(ch.dataLen)
		switch ch.fieldKeyword {
		case crtrFldTitle:
			d.creator.title = d.readString(int(ch.dataLen))
		case crtrFldCrtDate:
			d.creator.creationDate = time.Unix(int64(d.readUint32()), 0)
		case crtrFldModDate:
			d.creator.modificationDate = time.Unix(int64(d.readUint32()), 0)
		case crtrFldArtist:
			d.creator.artist = d.readString(int(ch.dataLen))
		case crtrFldCpyrght:
			d.creator.copyright = d.readString(int(ch.dataLen))
		case crtrFldDesc:
			d.creator.description = d.readString(int(ch.dataLen))
		case crtrFldAppID:
			d.creator.appID = d.readUint32()
		case crtrFldAppVer:
			d.creator.appVersion = d.readUint32()
		default:
			d.skip(int(ch.dataLen))
		}
	}
}

func (d *decoder) skip(n int) {
	_, err := d.r.Discard(n)
	if err != nil {
		d.error(err)
	}
}

func (d *decoder) read(b []byte) {
	if _, err := io.ReadFull(d.r, b); err != nil {
		d.error(err)
	}
}

func (d *decoder) readRect() image.Rectangle {
	d.read(d.tmpBuf[:16])
	return image.Rect(
		int(int32(decodeUint32(d.tmpBuf[:4]))),
		int(int32(decodeUint32(d.tmpBuf[4:8]))),
		int(int32(decodeUint32(d.tmpBuf[8:12]))),
		int(int32(decodeUint32(d.tmpBuf[12:16]))),
	)
}

func (d *decoder) readString(n int) string {
	// sanity check
	if n > 1024 {
		d.error(FormatError("bad string length"))
	}
	if cap(d.tmpBuf) < n {
		d.tmpBuf = make([]byte, n)
	}
	d.read(d.tmpBuf[:n])
	return string(d.tmpBuf[:n])
}

func (d *decoder) readByte() byte {
	b, err := d.r.ReadByte()
	if err != nil {
		d.error(err)
	}
	return b
}

func (d *decoder) readUint16() uint16 {
	d.read(d.tmpBuf[:2])
	return decodeUint16(d.tmpBuf[:2])
}

func (d *decoder) readUint32() uint32 {
	d.read(d.tmpBuf[:4])
	return decodeUint32(d.tmpBuf[:4])
}

func (d *decoder) readChunkHeader(ch *chunkHeader) {
	d.read(d.tmpBuf[:10])
	d.decodeChunkHeader(d.tmpBuf[:10], ch)
}

func (d *decoder) decodeChunkHeader(buf []byte, ch *chunkHeader) {
	if !bytes.Equal(buf[:4], chunkMagic) {
		d.error(FormatError("bad chunk magic"))
	}
	ch.fieldKeyword = decodeUint16(buf[4:6])
	ch.dataLen = decodeUint32(buf[6:10])
	fmt.Printf("CHUNK %+v\n", ch)
}

// readBlockHeader reads the next block from the file. it accepts a block
// rather than returning one so that the buffer can be reused.
func (d *decoder) readBlockHeader(bh *blockHeader) {
	if d.versionMajor > 3 {
		d.read(d.tmpBuf[:10])
		bh.initLen = 0xDEADBEEF
		bh.dataLen = decodeUint32(d.tmpBuf[6:10])
	} else {
		d.read(d.tmpBuf[:14])
		bh.initLen = decodeUint32(d.tmpBuf[6:10])
		bh.dataLen = decodeUint32(d.tmpBuf[10:14])
	}
	if !bytes.Equal(d.tmpBuf[:4], blockMagic) {
		d.error(FormatError("bad block magic"))
	}
	bh.id = blockID(decodeUint16(d.tmpBuf[4:6]))
	fmt.Printf("BLOCK %s %+v\n", bh.id, bh)
}

func decodeUint16(b []byte) uint16 {
	return uint16(b[0]) | (uint16(b[1]) << 8)
}

func decodeUint32(b []byte) uint32 {
	return uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}

func decodeUint64(b []byte) uint64 {
	return uint64(b[0]) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) | (uint64(b[5]) << 40) | (uint64(b[6]) << 48) | (uint64(b[7]) << 56)
}
