package confluence

const sngErrorsTpl = `
<p><div class="markdown-macro conf-macro output-inline" data-hasbody="true" data-macro-name="markdown" id="markdown-macro-1"><strong>此文档由服务构建时自动生成，请勿手动修改</strong><table class="confluenceTable">
 <thead>
  <tr>
   <th class="confluenceTh" style="text-align: center;">code</th>
   <th class="confluenceTh" style="text-align: center;">name</th>
   <th class="confluenceTh" style="text-align: center;">msg</th>
  </tr>
 </thead>
 <tbody>
{{- range $error := .errors }}   
   <tr>
	<td class="confluenceTd" style="text-align: center;">{{$error.Code}}</td>
	<td class="confluenceTd" style="text-align: center;">{{$error.Name}}</td>
	<td class="confluenceTd" style="text-align: center;">
      {{- range $msg := $error.Msgs }} 
      <p>{{$msg}}</p>
      {{- end }}
   </td>
   </tr>
{{- end }}
 </tbody>
</table></div></p>
`

const sngTestTpl = `
<p><strong>此文档由服务构建时自动生成，请勿手动修改</strong></p>
<div class="table-wrap">
    <table class="confluenceTable">
        <tbody>
            <tr>
                <th style="text-align: center;" class="confluenceTh">service name</th>
                <th style="text-align: center;" class="confluenceTh">version</th>
                <th style="text-align: center;" class="confluenceTh">unit test</th>
                <th style="text-align: center;" class="confluenceTh">api test</th>
                <th style="text-align: center;" class="confluenceTh">connserver test</th>
                <th style="text-align: center;" class="confluenceTh">gameserver test</th>
                <th style="text-align: center;" class="confluenceTh">benchmark test</th>
            </tr>
         {{- range $test := .tests }}   
            <tr>
               <td colspan="1" style="text-align: center;" class="confluenceTd">{{$test.Name}}</td>
               <td colspan="1" style="text-align: center;" class="confluenceTd">{{$test.Version}}</td>
               <td colspan="1" style="text-align: center;" class="confluenceTd"><a class="external-link"
                  href="{{$test.UnitTestReportUrl}}"
                  rel="nofollow">{{$test.Coverage}}</a></td>
               <td colspan="1" style="text-align: center;" class="confluenceTd">{{$test.HasApiTest}}</td>
               <td colspan="1" style="text-align: center;" class="confluenceTd">{{$test.HasConnserverTest}}</td>
               <td colspan="1" style="text-align: center;" class="confluenceTd">{{$test.HasGameserverTest}}</td>
               <td colspan="1" style="text-align: center;" class="confluenceTd"><a class="external-link"
                     href="{{$test.BenchmarkReportUrl}}"
                     rel="nofollow">{{$test.HasBenchmark}}</a></td>
            </tr>
         {{- end }}
        </tbody>
    </table>
</div>
`

const codeTpl = `
package code
var Codes = map[int]string{
   {{- range $code := .codes}} 
      {{$code.Code}}:"{{index $code.Msgs 0}}",

   {{- end}}
}
`

const codeJsonTpl = `
[
   {{- range $code := .codes}} 
   {
      "id": "{{$code.Code}}",
      "translation": "{{index $code.Msgs 0}}"
   },
   {{- end}}
]
`

const sngChangeLogTpl = `
<p><div class="markdown-macro conf-macro output-inline" data-hasbody="true" data-macro-name="markdown" id="markdown-macro-1"><strong>此文档由服务构建时自动生成，请勿手动修改</strong><table class="confluenceTable">
 <thead>
  <tr>
   <th class="confluenceTh" style="text-align: center;">version</th>
   <th class="confluenceTh" style="text-align: center;">releaseDate</th>
   <th class="confluenceTh" style="text-align: center;">msg</th>
  </tr>
 </thead>
 <tbody>
{{- range $log := .changelogs }}   
   <tr>
	<td class="confluenceTd" style="text-align: center;">{{$log.Version}}</td>
	<td class="confluenceTd" style="text-align: center;">{{$log.ReleaseDate}}</td>
	<td class="confluenceTd" style="text-align: center;">
      {{- range $msg := $log.Msgs }} 
      <p>{{$msg}}</p>
      {{- end }}
   </td>
   </tr>
{{- end }}
 </tbody>
</table></div></p>
`

const markdownChangeLogTpl = `
<ac:layout><ac:layout-section ac:type="two_equal"><ac:layout-cell>
<p class="auto-cursor-target"><br /></p><ac:structured-macro ac:name="tip" ac:schema-version="1" ac:macro-id="89af6ef0-0ce6-4179-8c04-570a19660e2b"><ac:rich-text-body>
<h4><strong>最新版本更新说明</strong></h4></ac:rich-text-body></ac:structured-macro>
<p>版本名称：{{ .LatestVersion}}</p>
<p>更新时间：{{.LatestReleaseDate}}</p>
<p>下载链接：<a href="{{.LatestDownloadLink}}" target="_blank">{{.LatestDownloadLink}}</a></p>
<p>更新内容：</p>
<ul>
{{- range $msg := .LatestMsgs }}
<li>{{$msg}}</li>
{{- end }}
</ul>
</ac:layout-cell><ac:layout-cell>
<p class="auto-cursor-target"><br /></p><ac:structured-macro ac:name="info" ac:schema-version="1" ac:macro-id="a45b3f43-812e-4d35-bc1b-3bef6fa27adb"><ac:rich-text-body>
<h4><strong>历史更新版本</strong></h4></ac:rich-text-body></ac:structured-macro>
{{- range $log := .HistoricalVersion }}
<p class="auto-cursor-target"><br /></p><ac:structured-macro ac:name="expand" ac:schema-version="1" ac:macro-id="73d2a9e4-7fce-4d3a-8b9b-bfbf574fb05c"><ac:parameter ac:name="title">{{ $log.Title }}</ac:parameter><ac:rich-text-body>
<p>更新时间：{{$log.ReleaseDate}}</p>
<p>下载链接：</p>
<p>更新内容：</p>
<ul>
{{- range $message := $log.Msgs }}
<li>{{$message}}</li>
{{- end }}
</ul>
</ac:rich-text-body></ac:structured-macro>
{{- end }}
<p class="auto-cursor-target"><br /></p></ac:layout-cell></ac:layout-section></ac:layout>
`
