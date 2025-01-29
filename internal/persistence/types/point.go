package types

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Point struct {
	Lat decimal.Decimal
	Lng decimal.Decimal
}

func NewPoint(lat, lng decimal.Decimal) *Point {
	return &Point{Lat: lat, Lng: lng}
}

func (p *Point) GormDataType() string {
	return "geography(Point, 4326)"
}

func (p *Point) GormDBDataType(*gorm.DB, *schema.Field) string {
	return "geography(Point, 4326)"
}

func (p *Point) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_GeographyFromText(?)",
		Vars: []any{p.String()},
	}
}

func (p *Point) Scan(val any) error {
	log.Debug().Msg("scan point")

	if _, ok := val.(string); !ok {
		return fmt.Errorf("invalid types %T", val)
	}

	b, err := hex.DecodeString(val.(string))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	coords := make([]decimal.Decimal, 2)
	if err := binary.Read(r, byteOrder, coords); err != nil {
		return err
	}

	p.Lng, p.Lat = coords[0], coords[1]
	return nil
}

func (p *Point) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p *Point) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Lng, p.Lat)
}
