include Nanoc::Helpers::Blogging
include Nanoc::Helpers::LinkTo
include Nanoc::Helpers::Rendering
include Nanoc::Helpers::Tagging

def date_ymd(time)
  case time
  when String
    return Time.parse(time).strftime("%Y-%m-%d")
  else
    return time
  end
end

class Nanoc::Item
  def moot_path
    path = File.dirname(identifier)
    slug = File.basename(identifier)
    "#{path}:#{slug}"
  end

  def summarize
    content = rep_named(:default).compiled_content

    # If there's an <hX> or <hr> anywhere in the content, assume that everything
    # before it is an introduction, and use that as the summary.
    if content =~ /<h[1-6r]/ then
      return content.split(/<h[1-6r]/).first
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
