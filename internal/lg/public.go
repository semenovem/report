package lg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type base struct {
	requestIDExtractor func(ctx context.Context) string
	cli                bool
	hideTime           bool
	level              int8
	outErr             io.Writer
	outInf             io.Writer
	outDeb             io.Writer
}

type Lg struct {
	base     *base
	names    []string
	params   []string
	isNested bool
	tags     []string
	ctx      context.Context
}

func New() (*Lg, *Setter) {
	st := Setter{
		logger: &Lg{
			base: &base{
				outErr: os.Stdout,
				outInf: os.Stdout,
				outDeb: os.Stdout,
			},
		},
	}

	return st.logger, &st
}

func (p *Lg) Func(ctx context.Context, n string) *Lg {
	a := p.copy()
	a.ctx = ctx

	if a.names == nil {
		a.names = make([]string, 0)
	}

	a.names = append(a.names, n)

	return a
}

func (p *Lg) Named(n string) *Lg {
	a := p.copy()

	if a.names == nil {
		a.names = make([]string, 0)
	}

	a.names = append(a.names, n)

	return a
}

func (p *Lg) With(k string, v interface{}) *Lg {
	a := p.copy()

	if b, err := json.Marshal(v); err != nil {
		a.with(k, fmt.Sprintf("%+v", v))
		p.save(Error, "inside error=[%s]", err.Error())
	} else {
		a.with(k, fmt.Sprintf("%+v", string(b)))
	}

	return a
}

func (p *Lg) Error(format string) {
	p.save(Error, format)
}

func (p *Lg) ErrorE(err error) {
	p.save(Error, err.Error())
}

func (p *Lg) Errorf(format string, v ...any) {
	p.save(Error, format, v...)
}

func (p *Lg) Info(format string) {
	p.save(Info, format)
}

func (p *Lg) Infof(format string, v ...any) {
	p.save(Info, format, v...)
}

func (p *Lg) Debug(format string) {
	p.save(Debug, format)
}

func (p *Lg) Debugf(format string, v ...any) {
	p.save(Debug, format, v...)
}

func (p *Lg) DebugOrErr(isDebug bool, format string) {
	if isDebug {
		p.save(Debug, format)
	} else {
		p.save(Error, format)
	}
}

func (p *Lg) DebugOrErrf(isDebug bool, format string, v ...any) {
	if isDebug {
		p.save(Debug, format, v...)
	} else {
		p.save(Error, format, v...)
	}
}

func (p *Lg) Nested(err error) {
	p.isNested = true
	p.save(Debug, err.Error())
}

func (p *Lg) NestedWith(err error, msg string) error {
	p.isNested = true

	if err == nil && msg == "" {
		p.save(Error, "logger argument is Nil (cause - developer)")
		return nil
	}

	if msg == "" {
		p.save(Debug, err.Error())
	} else {
		if err == nil {
			err = errors.New(msg)
		} else {
			p.with(nesterErrMsg, " "+err.Error())
			err = fmt.Errorf("%s: %s", err.Error(), msg)
		}

		p.save(Debug, msg)
	}

	return err
}

func (p *Lg) Nestedf(format string, v ...any) {
	p.isNested = true
	p.save(Debug, format, v...)
}

func (p *Lg) addTag(tags ...string) *Lg {
	if p.tags == nil {
		p.tags = make([]string, 0)
	}

	p.tags = append(p.tags, tags...)

	return p
}

// ----------------------------------------

func (p *Lg) DB(err error) {
	p.copy().addTag(databaseTag).save(Error, err.Error())
}

func (p *Lg) DBStr(msg string) {
	p.copy().addTag(databaseTag).save(Error, msg)
}

func (p *Lg) DBf(format string, v ...any) {
	p.copy().addTag(databaseTag).save(Error, format, v...)
}

func (p *Lg) Redis(err error) {
	p.copy().addTag(redisTag).save(Error, err.Error())
}

func (p *Lg) RedisStr(msg string) {
	p.copy().addTag(redisTag).save(Error, msg)
}

func (p *Lg) Redisf(format string, v ...any) {
	p.copy().addTag(redisTag).save(Error, format, v...)
}

func (p *Lg) BadRequest(err error) {
	p.copy().addTag(badRequestTag).save(Debug, err.Error())
}

func (p *Lg) BadRequestStr(msg string) {
	p.copy().addTag(badRequestTag).save(Debug, msg)
}

func (p *Lg) BadRequestStrRetErr(msg string) error {
	p.copy().addTag(badRequestTag).BadRequestStr(msg)
	return errors.New(msg)
}

func (p *Lg) NotFound(err error) {
	p.copy().addTag(notFound).save(Info, err.Error())
}

func (p *Lg) NotFoundStr(msg string) {
	p.copy().addTag(notFound).save(Info, msg)
}

func (p *Lg) Deny(err error) {
	p.copy().addTag(denyTag).save(Info, err.Error())
}

func (p *Lg) Auth(err error) {
	p.copy().addTag(authTag).save(Info, err.Error())
}

func (p *Lg) AuthStr(msg string) {
	p.copy().addTag(authTag).save(Info, msg)
}

func (p *Lg) AuthDebug(err error) {
	p.copy().addTag(authTag).save(Debug, err.Error())
}

func (p *Lg) AuthDebugStr(msg string) {
	p.copy().addTag(authTag).save(Debug, msg)
}
