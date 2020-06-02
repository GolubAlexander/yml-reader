package ymlreader

import (
	"encoding/xml"
	"time"
)

const datePattern = "2006-01-02 15:04"

type YMLCatalog struct {
	Date YMLDate `xml:"date,attr"`
	Shop Shop    `xml:"shop"`
}

type Shop struct {
	Name            string           `xml:"name"`
	Company         string           `xml:"company"`
	URL             string           `xml:"url"`
	Currencies      []Currency       `xml:"currencies>currency"`
	Categories      []Category       `xml:"categories>category"`
	DeliveryOptions []DeliveryOption `xml:"delivery-options>option"`
	Offers          []Offer          `xml:"offers>offer"`
}

type Offer struct {
	Type                 string           `xml:"type,attr"`
	Available            string           `xml:"available,attr"`
	AttrID               string           `xml:"id,attr"`
	GroupID              string           `xml:"group_id,attr"`
	URL                  string           `xml:"url"`
	Price                string           `xml:"price"`
	CurrencyId           string           `xml:"currencyId"`
	CategoryId           string           `xml:"categoryId"`
	Picture              string           `xml:"picture"`
	Store                string           `xml:"store"`
	Pickup               string           `xml:"pickup"`
	Delivery             string           `xml:"delivery"`
	Author               string           `xml:"author"`
	Name                 string           `xml:"name"`
	ID                   string           `xml:"ID"`
	Publisher            string           `xml:"publisher"`
	Series               string           `xml:"series"`
	Year                 string           `xml:"year"`
	Description          string           `xml:"description"`
	ISBN                 string           `xml:"ISBN"`
	Language             string           `xml:"language"`
	Binding              string           `xml:"binding"`
	PageExtent           string           `xml:"page_extent"`
	ManufacturerWarranty string           `xml:"manufacturer_warranty"`
	Barcode              string           `xml:"barcode"`
	Weight               string           `xml:"weight"`
	DeliveryOptions      []DeliveryOption `xml:"delivery-options>option"`
	Param                struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"param"`
	Dimensions string `xml:"dimensions"`
}

type Currency struct {
	ID   string `xml:"id,attr"`
	Rate int    `xml:"rate,attr"`
}

type Category struct {
	Name     string `xml:",chardata"`
	ID       int    `xml:"id,attr"`
	ParentID int    `xml:"parentId,attr"`
}

type DeliveryOption struct {
	Cost        string `xml:"cost,attr"`
	Days        string `xml:"days,attr"`
	OrderBefore string `xml:"order-before,attr"`
}

type YMLDate struct {
	time.Time
}

func (y *YMLDate) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, _ := time.Parse(datePattern, attr.Value)
	*y = YMLDate{parse}
	return nil
}
