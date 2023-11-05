package lg

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"
)

func (p *Lg) save(level Level, format string, v ...any) {
	var out = make([]byte, prefixLen)

	if !p.base.hideTime {
		out = append(out, []byte(fmt.Sprintf(" [%s]", time.Now().Format(time.RFC3339)))...)
	}

	if p.ctx != nil && p.base.requestIDExtractor != nil {
		s := " [req_id:" + p.base.requestIDExtractor(p.ctx) + "]"
		out = append(out, []byte(s)...)
	}

	// ------------------------------------------------------
	// Оформление caller
	_, path, line, _ := runtime.Caller(2)
	ps := strings.Split(path, "/")

	if len(ps) > 2 {
		path = strings.Join(ps[len(ps)-3:], "/")
	}

	caller := fmt.Sprintf(" [%s:%d] ", path, line)
	out = append(out, []byte(caller)...)
	// ------------------------------------------------------

	if p.names != nil {
		out = append(out, []byte(" "+strings.Join(p.names, "."))...)
		if p.isNested {
			out = append(out, nestedBytes...)
		}
	}

	if p.tags != nil {
		out = append(out, []byte(" ["+strings.Join(p.tags, ".")+"]")...)
		if p.isNested {
			out = append(out, nestedBytes...)
		}
	}

	if p.params != nil {
		out = append(out, []byte(" ["+strings.Join(p.params, ", ")+"]")...)
	}

	out = append(out, []byte(" "+fmt.Sprintf(format, v...))...)

	out = append(out, []byte("\n")...)

	var (
		writer io.Writer
		prefix []byte
	)

	switch level {
	case Error:
		prefix = prefixBytesErr
		writer = p.base.outErr
	case Info:
		prefix = prefixBytesInfo
		writer = p.base.outInf
	case Debug:
		prefix = prefixBytesDebug
		writer = p.base.outDeb
	}

	_ = copy(out[0:prefixLen], prefix)

	if _, err := writer.Write(out); err != nil {
		fmt.Println(prefixErr, string(out))
	}
}

func (p *Lg) copy() *Lg {
	a := Lg{
		base:     p.base,
		isNested: p.isNested,
		ctx:      p.ctx,
		// TODO нужно скопировать оставшиеся поля
	}

	if p.names != nil {
		a.names = make([]string, len(p.names))
		copy(a.names, p.names)
	}

	if p.params != nil {
		a.params = make([]string, len(p.params))
		copy(a.params, p.params)
	}

	if p.tags != nil {
		a.tags = make([]string, len(p.tags))
		copy(a.tags, p.tags)
	}

	return &a
}

func (p *Lg) with(k, v string) {
	if p.params == nil {
		p.params = make([]string, 0)
	}

	p.params = append(p.params, fmt.Sprintf("%s:%s", k, v))
}

func (p *Lg) nestedArg(with any) {
	var s string

	switch t := with.(type) {
	case string:
		s = t
	case error:
		s = t.Error()
	default:
		s = fmt.Sprintf("%+v", with)
	}

	p.with(nesterErrMsg, s)
}
