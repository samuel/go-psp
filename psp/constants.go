package psp

import "fmt"

// Block identifiers (PSPBlockID)
type blockID uint16

const (
	imageBlock               blockID = iota // General Image Attributes Block (main)
	creatorBlock                            // Creator Data Block (main)
	colorBlock                              // Color Palette Block (main and sub)
	layerStartBlock                         // Layer Bank Block (main)
	layerBlock                              // Layout Block (sub)
	channelBlock                            // Channel block (sub)
	selectionBlock                          // Selection block (main)
	alphaBankBlock                          // Alpha bank block (main)
	alphaChannelBlock                       // Alpha Channel Block (sub)
	thumbnailBlock                          // Thumbnail Block (main)
	extendedDataBlock                       // Extended Data Block (main)
	tubeBlock                               // Picture Tube Data Block (main)
	adjustmentExtensionBlock                // Adjustment Layer Extension Block (sub) (since PSP6
	vectorExtensionBlock                    // Vector Layer Extension Block (sub) (since PSP6)
	shapeBlock                              // Vector Shape Block (sub) (since PSP6)
	paintstyleBlock                         // Paint Style Block (sub) (since PSP6)
	compositeImageBankBlock                 // Composite Image Bank (main) (since PSP6)
	compositeAttributesBlock                // Composite Image Attributes (sub) (since PSP6)
	jpegBlock                               // JPEG Image Block (sub) (since PSP6)
	linestyleBlock                          // Line Style Block (sub) (since PSP7)
	tableBankBlock                          // Table Bank Block (main) (since PSP7)
	tableBlock                              // Table Block (sub) (since PSP7)
	paperBlock                              // Vector Table Paper Block (sub) (since PSP7)
	patternBlock                            // Vector Table Pattern Block (sub) (since PSP7)
	gradientBlock                           // Vector Table Gradient Block (not used) (since PSP8)
	groupExtensionBlock                     // Group Layer Block (sub) (since PSP8)
	maskExtensionBlock                      // Mask Layer Block (sub) (since PSP8)
	brushBlock                              // Brush Data Block (main) (since PSP8)
)

var blockTypes = map[blockID]string{
	imageBlock:               "imageBlock",
	creatorBlock:             "creatorBlock",
	colorBlock:               "colorBlock",
	layerStartBlock:          "layerStartBlock",
	layerBlock:               "layerBlock",
	channelBlock:             "channelBlock",
	selectionBlock:           "selectionBlock",
	alphaBankBlock:           "alphaBankBlock",
	alphaChannelBlock:        "alphaChannelBlock",
	thumbnailBlock:           "thumbnailBlock",
	extendedDataBlock:        "extendedDataBlock",
	tubeBlock:                "tubeBlock",
	adjustmentExtensionBlock: "adjustmentExtensionBlock",
	vectorExtensionBlock:     "vectorExtensionBlock",
	shapeBlock:               "shapeBlock",
	paintstyleBlock:          "paintstyleBlock",
	compositeImageBankBlock:  "compositeImageBankBlock",
	compositeAttributesBlock: "compositeAttributesBlock",
	jpegBlock:                "jpegBlock",
	linestyleBlock:           "linestyleBlock",
	tableBankBlock:           "tableBankBlock",
	tableBlock:               "tableBlock",
	paperBlock:               "paperBlock",
	patternBlock:             "patternBlock",
	gradientBlock:            "gradientBlock",
	groupExtensionBlock:      "groupExtensionBlock",
	maskExtensionBlock:       "maskExtensionBlock",
	brushBlock:               "brushBlock",
}

func (id blockID) String() string {
	if s := blockTypes[id]; s != "" {
		return s
	}
	return fmt.Sprintf("blockID(%d)", id)
}

// Bitmap type (PSPDIBType)
type bitmapType uint16

const (
	dibImage              bitmapType = iota // Layer color bitmap
	dibTransMask                            // Layer transparency mask bitmap
	dibUserMask                             // Layer user mask bitmap
	dibSelection                            // Selection mask bitmap
	dibAlphaMask                            // Alpha channel mask bitmap
	dibThumbnail                            // Thumbnail bitmap
	dibThumbnailTransMask                   // Thumbnail transparency mask (since PSP6)
	dibAdjustmentLayer                      // Adjustment layer bitmap (since PSP6)
	dibComposite                            // Composite image bitmap (since PSP6)
	dibCompositeTransMask                   // Composite image transparency (since PSP6)
	dibPaper                                // Paper bitmap (since PSP7)
	dibPattern                              // Pattern bitmap (since PSP7)
	dibPatternTransMask                     // Pattern transparency mask (since PSP7)
)

var bitmapTypes = map[bitmapType]string{
	dibImage:              "dibImage",
	dibTransMask:          "dibTransMask",
	dibUserMask:           "dibUserMAsk",
	dibSelection:          "dibSelection",
	dibAlphaMask:          "dibAlphaMask",
	dibThumbnail:          "dibThumbnail",
	dibThumbnailTransMask: "dibThumbnailTransMask",
	dibAdjustmentLayer:    "dibAdjustmentLayer",
	dibComposite:          "dibComposite",
	dibCompositeTransMask: "dibCompositeTransMask",
	dibPaper:              "dibPaper",
	dibPattern:            "dibPattern",
	dibPatternTransMask:   "dibPatternTransMask",
}

func (bt bitmapType) String() string {
	if s := bitmapTypes[bt]; s != "" {
		return s
	}
	return fmt.Sprintf("bitmapType(%d)", bt)
}

// Channel types (PSPChannelType)
type channelType uint16

const (
	channelComposite channelType = iota // Channel of single channel bitmap
	channelRed                          // Red channel of 24 bit bitmap
	channelGreen                        // Green channel of 24 bit bitmap
	channelBlue                         // Blue channel of 24 bit bitmap
)

func (ct channelType) String() string {
	switch ct {
	case channelComposite:
		return "channelComposite"
	case channelRed:
		return "channelRed"
	case channelGreen:
		return "channelGreen"
	case channelBlue:
		return "channelBlue"
	}
	return fmt.Sprintf("channelType(%d)", ct)
}

// Possible metrics used to measure resolution. (PSP_METRIC)
type metric byte

const (
	metricUndefined metric = iota
	metricInch
	metricCentimeters
)

// Possible types of compression (PSPCompression)
type compression uint16

const (
	compressionNone compression = iota
	compressionRLE
	compressionLZ77
)

// Picture tube placement mode (TubePlacementMode)
const (
	tpmRandom   = iota // Place tube images in random intervals
	tpmConstant        // Place tube images in constant intervals
)

// Picture tube selection mode (TubeSelectionMode)
const (
	tsmRandom      = iota // Randomly select the next image in tube to display
	tsmIncremental        // Select each tube image in turn
	tsmAngular            // Select image based on cursor direction
	tsmPressure           // Select image based on pressure (from pressure-sensitive pad)
	tsmVelocity           // Select image based on cursor speed/* Extended data field types.
)

// Extended data field types (PSPExtendedDataID)
const (
	xDataTrnsIndex = iota // Transparency index field
)

// Creator field types (PSPCreatorFieldID)
const (
	crtrFldTitle   = iota // Image document title field
	crtrFldCrtDate        // Creation date field
	crtrFldModDate        // Modification date field
	crtrFldArtist         // Artist name field
	crtrFldCpyrght        // Copyright holder name field
	crtrFldDesc           // Image document description field
	crtrFldAppID          // Creating app id field
	crtrFldAppVer         // Creating app version field
)

// Creator application identifiers (PSPCreatorAppID)
const (
	creatorAppUnknown      = iota // Creator application unknown
	creatorAppPaintShopPro        // Creator is Paint Shop Pro
)

// Layer types (PSPLayerType)
type layerType byte

const (
	layerNormal            layerType = iota // Normal layer
	layerFloatingSelection                  // Floating selection layer
)

func (lt layerType) String() string {
	switch lt {
	case layerNormal:
		return "layerNormal"
	case layerFloatingSelection:
		return "layerFloatingSelection"
	}
	return fmt.Sprintf("layerType(%d)", lt)
}

// /* Graphic contents flags. (since PSP6)
//  */
// typedef enum {
//   /* Layer types */
//   keGCRasterLayers     = 0x00000001,    /* At least one raster layer */
//   keGCVectorLayers     = 0x00000002,    /* At least one vector layer */
//   keGCAdjustmentLayers = 0x00000004,    /* At least one adjustment layer */

//   /* Additional attributes */
//   keGCThumbnail              = 0x01000000,      /* Has a thumbnail */
//   keGCThumbnailTransparency  = 0x02000000,      /* Thumbnail transp. */
//   keGCComposite              = 0x04000000,      /* Has a composite image */
//   keGCCompositeTransparency  = 0x08000000,      /* Composite transp. */
//   keGCFlatImage              = 0x10000000,      /* Just a background */
//   keGCSelection              = 0x20000000,      /* Has a selection */
//   keGCFloatingSelectionLayer = 0x40000000,      /* Has float. selection */
//   keGCAlphaChannels          = 0x80000000,      /* Has alpha channel(s) */
// } PSPGraphicContents;

// /* Character style flags. (since PSP6)
//  */
// typedef enum {
//   keStyleItalic      = 0x00000001,      /* Italic property bit */
//   keStyleStruck      = 0x00000002,      /* Strike­out property bit */
//   keStyleUnderlined  = 0x00000004,      /* Underlined property bit */
//   keStyleWarped      = 0x00000008,      /* Warped property bit (since PSP8) */
//   keStyleAntiAliased = 0x00000010,      /* Anti­aliased property bit (since PSP8) */
// } PSPCharacterProperties;

// /* Table type. (since PSP7)
//  */
// typedef enum {
//   keTTUndefined = 0,     /* Undefined table type */
//   keTTGradientTable,     /* Gradient table type */
//   keTTPaperTable,        /* Paper table type */
//   keTTPatternTable       /* Pattern table type */
// } PSPTableType;

// /* Layer flags. (since PSP6)
//  */
// typedef enum {
//   keVisibleFlag      = 0x00000001,      /* Layer is visible */
//   keMaskPresenceFlag = 0x00000002,      /* Layer has a mask */
// } PSPLayerProperties;

// /* Shape property flags. (since PSP6)
//  */
// typedef enum {
//   keShapeAntiAliased = 0x00000001,      /* Shape is anti­aliased */
//   keShapeSelected    = 0x00000002,      /* Shape is selected */
//   keShapeVisible     = 0x00000004,      /* Shape is visible */
// } PSPShapeProperties;

// /* Polyline node type flags. (since PSP7)
//  */
// typedef enum {
//   keNodeUnconstrained     = 0x0000,     /* Default node type */
//   keNodeSmooth            = 0x0001,     /* Node is smooth */
//   keNodeSymmetric         = 0x0002,     /* Node is symmetric */
//   keNodeAligned           = 0x0004,     /* Node is aligned */
//   keNodeActive            = 0x0008,     /* Node is active */
//   keNodeLocked            = 0x0010,     /* Node is locked */
//   keNodeSelected          = 0x0020,     /* Node is selected */
//   keNodeVisible           = 0x0040,     /* Node is visible */
//   keNodeClosed            = 0x0080,     /* Node is closed */

//   /* TODO: This might be a thinko in the spec document only or in the image
//    *       format itself. Need to investigate that later
//    */
//   keNodeLockedPSP6        = 0x0016,     /* Node is locked */
//   keNodeSelectedPSP6      = 0x0032,     /* Node is selected */
//   keNodeVisiblePSP6       = 0x0064,     /* Node is visible */
//   keNodeClosedPSP6        = 0x0128,     /* Node is closed */

// } PSPPolylineNodeTypes;

// /* Blend modes. (since PSP6)
//  */
// typedef enum {
//   PSP_BLEND_NORMAL,
//   PSP_BLEND_DARKEN,
//   PSP_BLEND_LIGHTEN,
//   PSP_BLEND_HUE,
//   PSP_BLEND_SATURATION,
//   PSP_BLEND_COLOR,
//   PSP_BLEND_LUMINOSITY,
//   PSP_BLEND_MULTIPLY,
//   PSP_BLEND_SCREEN,
//   PSP_BLEND_DISSOLVE,
//   PSP_BLEND_OVERLAY,
//   PSP_BLEND_HARD_LIGHT,
//   PSP_BLEND_SOFT_LIGHT,
//   PSP_BLEND_DIFFERENCE,
//   PSP_BLEND_DODGE,
//   PSP_BLEND_BURN,
//   PSP_BLEND_EXCLUSION,
//   PSP_BLEND_TRUE_HUE, /* since PSP8 */
//   PSP_BLEND_TRUE_SATURATION, /* since PSP8 */
//   PSP_BLEND_TRUE_COLOR, /* since PSP8 */
//   PSP_BLEND_TRUE_LIGHTNESS, /* since PSP8 */
//   PSP_BLEND_ADJUST = 255,
// } PSPBlendModes;

// /* Adjustment layer types. (since PSP6)
//  */
// typedef enum {
//   keAdjNone = 0,        /* Undefined adjustment layer type */
//   keAdjLevel,           /* Level adjustment */
//   keAdjCurve,           /* Curve adjustment */
//   keAdjBrightContrast,  /* Brightness­contrast adjustment */
//   keAdjColorBal,        /* Color balance adjustment */
//   keAdjHSL,             /* HSL adjustment */
//   keAdjChannelMixer,    /* Channel mixer adjustment */
//   keAdjInvert,          /* Invert adjustment */
//   keAdjThreshold,       /* Threshold adjustment */
//   keAdjPoster           /* Posterize adjustment */
// } PSPAdjustmentLayerType;

// /* Vector shape types. (since PSP6)
//  */
// typedef enum {
//   keVSTUnknown = 0,     /* Undefined vector type */
//   keVSTText,            /* Shape represents lines of text */
//   keVSTPolyline,        /* Shape represents a multiple segment line */
//   keVSTEllipse,         /* Shape represents an ellipse (or circle) */
//   keVSTPolygon,         /* Shape represents a closed polygon */
//   keVSTGroup,           /* Shape represents a group shape (since PSP7) */
// } PSPVectorShapeType;

// /* Text element types. (since PSP6)
//  */
// typedef enum {
//   keTextElemUnknown = 0,        /* Undefined text element type */
//   keTextElemChar,               /* A single character code */
//   keTextElemCharStyle,          /* A character style change */
//   keTextElemLineStyle           /* A line style change */
// } PSPTextElementType;

// /* Text alignment types. (since PSP6)
//  */
// typedef enum {
//   keTextAlignmentLeft = 0,      /* Left text alignment */
//   keTextAlignmentCenter,        /* Center text alignment */
//   keTextAlignmentRight          /* Right text alignment */
// } PSPTextAlignment;

// /* Paint style types. (since PSP6)
//  */
// typedef enum {
//   keStyleNone     = 0x0000,     /* No paint style info applies */
//   keStyleColor    = 0x0001,     /* Color paint style info */
//   keStyleGradient = 0x0002,     /* Gradient paint style info */
//   keStylePattern  = 0x0004,     /* Pattern paint style info (since PSP7) */
//   keStylePaper    = 0x0008,      Paper paint style info (since PSP7)
//   keStylePen      = 0x0010,     /* Organic pen paint style info (since PSP7) */
// } PSPPaintStyleType;

// /* Gradient type. (since PSP7)
//  */
// typedef enum {
//   keSGTLinear = 0,      /* Linera gradient type */
//   keSGTRadial,          /* Radial gradient type */
//   keSGTRectangular,     /* Rectangulat gradient type */
//   keSGTSunburst         /* Sunburst gradient type */
// } PSPStyleGradientType;

// /* Paint Style Cap Type (Start & End). (since PSP7)
//  */
// typedef enum {
//   keSCTCapFlat = 0,             /* Flat cap type (was round in psp6) */
//   keSCTCapRound,                /* Round cap type (was square in psp6) */
//   keSCTCapSquare,               /* Square cap type (was flat in psp6) */
//   keSCTCapArrow,                /* Arrow cap type */
//   keSCTCapCadArrow,             /* Cad arrow cap type */
//   keSCTCapCurvedTipArrow,       /* Curved tip arrow cap type */
//   keSCTCapRingBaseArrow,        /* Ring base arrow cap type */
//   keSCTCapFluerDelis,           /* Fluer deLis cap type */
//   keSCTCapFootball,             /* Football cap type */
//   keSCTCapXr71Arrow,            /* Xr71 arrow cap type */
//   keSCTCapLilly,                /* Lilly cap type */
//   keSCTCapPinapple,             /* Pinapple cap type */
//   keSCTCapBall,                 /* Ball cap type */
//   keSCTCapTulip                 /* Tulip cap type */
// } PSPStyleCapType;

// /* Paint Style Join Type. (since PSP7)
//  */
// typedef enum {
//   keSJTJoinMiter = 0,
//   keSJTJoinRound,
//   keSJTJoinBevel
// } PSPStyleJoinType;

// /* Organic pen type. (since PSP7)
//  */
// typedef enum {
//   keSPTOrganicPenNone = 0,      /* Undefined pen type */
//   keSPTOrganicPenMesh,          /* Mesh pen type */
//   keSPTOrganicPenSand,          /* Sand pen type */
//   keSPTOrganicPenCurlicues,     /* Curlicues pen type */
//   keSPTOrganicPenRays,          /* Rays pen type */
//   keSPTOrganicPenRipple,        /* Ripple pen type */
//   keSPTOrganicPenWave,          /* Wave pen type */
//   keSPTOrganicPen               /* Generic pen type */
// } PSPStylePenType;

// /* Channel types.
//  */
// typedef enum {
//   PSP_CHANNEL_COMPOSITE = 0,    /* Channel of single channel bitmap */
//   PSP_CHANNEL_RED,              /* Red channel of 24 bit bitmap */
//   PSP_CHANNEL_GREEN,            /* Green channel of 24 bit bitmap */
//   PSP_CHANNEL_BLUE              /* Blue channel of 24 bit bitmap */
// } PSPChannelType;

// /* Possible metrics used to measure resolution.
//  */
// typedef enum {
//   PSP_METRIC_UNDEFINED = 0,     /* Metric unknown */
//   PSP_METRIC_INCH,              /* Resolution is in inches */
//   PSP_METRIC_CM                 /* Resolution is in centimeters */
// } PSP_METRIC;

// /* Possible types of compression.
//  */
// typedef enum {
//   PSP_COMP_NONE = 0,            /* No compression */
//   PSP_COMP_RLE,                 /* RLE compression */
//   PSP_COMP_LZ77,                /* LZ77 compression */
//   PSP_COMP_JPEG                 /* JPEG compression (only used by thumbnail and composite image) (since PSP6) */
// } PSPCompression;

// /* Picture tube placement mode.
//  */
// typedef enum {
//   tpmRandom,                    /* Place tube images in random intervals */
//   tpmConstant                   /* Place tube images in constant intervals */
// } TubePlacementMode;

// /* Picture tube selection mode.
//  */
// typedef enum {
//   tsmRandom,                    /* Randomly select the next image in  */
//                                 /* tube to display */
//   tsmIncremental,               /* Select each tube image in turn */
//   tsmAngular,                   /* Select image based on cursor direction */
//   tsmPressure,                  /* Select image based on pressure  */
//                                 /* (from pressure-sensitive pad) */
//   tsmVelocity                   /* Select image based on cursor speed */
// } TubeSelectionMode;

// /* Extended data field types.
//  */
// typedef enum {
//   PSP_XDATA_TRNS_INDEX = 0,     /* Transparency index field */
//   PSP_XDATA_GRID,               /* Image grid information (since PSP7) */
//   PSP_XDATA_GUIDE,              /* Image guide information (since PSP7) */
//   PSP_XDATA_EXIF,               /* Image Exif information (since PSP8) */
// } PSPExtendedDataID;

// /* Creator field types.
//  */
// typedef enum {
//   PSP_CRTR_FLD_TITLE = 0,       /* Image document title field */
//   PSP_CRTR_FLD_CRT_DATE,        /* Creation date field */
//   PSP_CRTR_FLD_MOD_DATE,        /* Modification date field */
//   PSP_CRTR_FLD_ARTIST,          /* Artist name field */
//   PSP_CRTR_FLD_CPYRGHT,         /* Copyright holder name field */
//   PSP_CRTR_FLD_DESC,            /* Image document description field */
//   PSP_CRTR_FLD_APP_ID,          /* Creating app id field */
//   PSP_CRTR_FLD_APP_VER          /* Creating app version field */
// } PSPCreatorFieldID;

// /* Grid units type. (since PSP7)
//  */
// typedef enum {
//   keGridUnitsPixels = 0,        /* Grid units is pixels */
//   keGridUnitsInches,            /* Grid units is inches */
//   keGridUnitsCentimeters        /* Grid units is centimeters */
// } PSPGridUnitsType;

// /* Guide orientation type. (since PSP7)
//  */
// typedef enum  {
//   keHorizontalGuide = 0,
//   keVerticalGuide
// } PSPGuideOrientationType;

// /* Creator application identifiers.
//  */
// typedef enum {
//   PSP_CREATOR_APP_UNKNOWN = 0,  /* Creator application unknown */
//   PSP_CREATOR_APP_PAINT_SHOP_PRO /* Creator is Paint Shop Pro */
// } PSPCreatorAppID;

// /* Layer types.
//  */
// typedef enum {
//   PSP_LAYER_NORMAL = 0,         /* Normal layer */
//   PSP_LAYER_FLOATING_SELECTION  /* Floating selection layer */
// } PSPLayerTypePSP5;

// /* Layer types. (since PSP6)
//  */
// typedef enum {
//   keGLTUndefined = 0,           /* Undefined layer type */
//   keGLTRaster,                  /* Standard raster layer */
//   keGLTFloatingRasterSelection, /* Floating selection (raster layer) */
//   keGLTVector,                  /* Vector layer */
//   keGLTAdjustment,              /* Adjustment layer */
//   keGLTMask                     /* Mask layer (since PSP8) */
// } PSPLayerTypePSP6;
