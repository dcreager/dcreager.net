require 'rocco'

class RoccoFilter < Nanoc::Filter
  identifier :rocco

  def run(content, params={})
    filename = @item.identifier.chop + '.' + @item[:extension]
    doc = Rocco.new(filename, [], params) do content.dup end
    doc.to_html
  end
end
