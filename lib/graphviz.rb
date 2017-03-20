class DotFilter < Nanoc::Filter
  identifier :dot
  type :binary

  def run(filename, params={})
    system(
      'dot',
      '-Tpng',
      '-Gdpi=250',
      filename, '-o', output_filename,
    )
  end
end
