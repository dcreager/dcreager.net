require 'coderay'

include Nanoc::Helpers::Capturing
include Nanoc::Helpers::Filtering
include Nanoc::Helpers::HTMLEscape

def highlight(lang, args={}, &block)
  erbout = eval('_erbout', block.binding)
  content = capture(&block).sub(/^\n/, "")
  colorized = ::CodeRay.scan(content, lang).html(args)
  if args[:line_numbers] then
    erbout << "<div class=\"CodeRay\">"
    erbout << colorized
    erbout << "</div>"
  else
    erbout << "<pre class=\"CodeRay\"><code>"
    erbout << colorized
    erbout << "</code></pre>"
  end
end
