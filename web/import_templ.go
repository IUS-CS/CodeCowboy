// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package web

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func (w *Web) importForm() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h1>Import Course Data</h1><form hx-post=\"/import\" enctype=\"multipart/form-data\"><div class=\"grid-x\"><div class=\"cell\">Course Name</div><div class=\"cell\"><input type=\"text\" name=\"course\"></div></div><div class=\"grid-x\"><div class=\"cell\">GitHub export</div><div class=\"cell\"><input type=\"file\" name=\"github\" accept=\"text/csv\"></div></div><div class=\"grid-x\"><div class=\"cell\">Canvas export</div><div class=\"cell\"><input type=\"file\" name=\"canvas\" accept=\"text/csv\"></div></div><div class=\"grid-x\"><div class=\"cell\">Assignments JSON</div><div class=\"cell\"><input type=\"file\" name=\"assignments\" accept=\"application/json\"></div></div><button type=\"submit\" class=\"submit button\">Save</button></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
