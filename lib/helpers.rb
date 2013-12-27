include Nanoc::Helpers::Blogging
include Nanoc::Helpers::LinkTo
include Nanoc::Helpers::Rendering
include Nanoc::Helpers::Tagging

class Nanoc::Item
  def summarize
    content = rep_named(:default).compiled_content

    # If there's an <h2> anywhere in the content, assume that everything before
    # it is an introduction, and use that as the summary.
    if content =~ /<h2/ then
      return content.split("<h2").first
    end

    # Otherwise, use the article's first four paragraphs as the summary.
    content.split("\n\n").take(4).join "\n\n"
  end
end

def navbar_link(text, target, attributes={})
  path = target.is_a?(String) ? target : target.path

  if @item_rep && @item_rep.path == path
    klass = " class=\"active\""
  else
    klass = ""
  end

  link = link_to(text, target, attributes)
  "<li#{klass}>#{link}</li>"
end
