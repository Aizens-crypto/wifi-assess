package wifi

import "github.com/google/gopacket/layers"

// FrameType is a thin, project-specific wrapper around gopacket's Dot11Type
// so the rest of wifi-assess depends on our own type rather than reaching
// into gopacket/layers directly everywhere.
type FrameType layers.Dot11Type

// Category groups related frame types together. 802.11 defines three.
type Category int

const (
	CategoryUnknown Category = iota
	CategoryManagement
	CategoryControl
	CategoryData
)

// The subset of management frame types wifi-assess currently cares about.
// Extend as new rules/detections need more types. Data frame subtypes are
// numerous (Data, QoS Data, Null, CF-Ack, ...) and aren't individually
// named here yet — Category() is enough to bucket them as "Data" for now.
const (
	FrameTypeBeacon             FrameType = FrameType(layers.Dot11TypeMgmtBeacon)
	FrameTypeProbeRequest       FrameType = FrameType(layers.Dot11TypeMgmtProbeReq)
	FrameTypeProbeResponse      FrameType = FrameType(layers.Dot11TypeMgmtProbeResp)
	FrameTypeAuthentication     FrameType = FrameType(layers.Dot11TypeMgmtAuthentication)
	FrameTypeDeauthentication   FrameType = FrameType(layers.Dot11TypeMgmtDeauthentication)
	FrameTypeAssociationRequest FrameType = FrameType(layers.Dot11TypeMgmtAssociationReq)
	FrameTypeAssociationResp    FrameType = FrameType(layers.Dot11TypeMgmtAssociationResp)
	FrameTypeDisassociation     FrameType = FrameType(layers.Dot11TypeMgmtDisassociation)
)

// Category classifies a frame type into management/control/data using
// gopacket's own MainType(), so we don't duplicate the 802.11 spec table.
func (f FrameType) Category() Category {
	switch layers.Dot11Type(f).MainType() {
	case layers.Dot11TypeMgmt:
		return CategoryManagement
	case layers.Dot11TypeCtrl:
		return CategoryControl
	case layers.Dot11TypeData:
		return CategoryData
	default:
		return CategoryUnknown
	}
}

// String delegates to gopacket's own frame type names so this stays in
// sync with the 802.11 spec table gopacket already maintains.
func (f FrameType) String() string {
	return layers.Dot11Type(f).String()
}

func (c Category) String() string {
	switch c {
	case CategoryManagement:
		return "Management"
	case CategoryControl:
		return "Control"
	case CategoryData:
		return "Data"
	default:
		return "Unknown"
	}
}
