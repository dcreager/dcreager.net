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

def full_title(item)
  if item[:series] then
    item[:series] + ': ' + item[:title]
  else
    item[:title]
  end
end

def summarize(item)
  content = item.compiled_content(snapshot: :pre)

  # If there's an <hX> or <hr> anywhere in the content, assume that everything
  # before it is an introduction, and use that as the summary.
  if content =~ /<h[1-6r]/ then
    return content.split(/<h[1-6r]/).first
  end

  # Otherwise, use the article's first four paragraphs as the summary.
  content.split("\n\n").take(4).join "\n\n"
end

def figure(slug)
  path = '/media/images/' + slug
  [
    "<picture>",
    "<source media=\"(max-width: 511px)\" srcset=\"#{path}.tall.png\">",
    "<source media=\"(min-width: 2048pw)\" srcset=\"#{path}.large.png\">",
    "<img src=\"#{path}.png\" class=\"figure\">",
    "</picture>",
  ].join('')
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
