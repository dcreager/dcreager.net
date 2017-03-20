class DotFilter < Nanoc::Filter
  identifier :dot
  type :binary

  SIZES = {
    :tall => ['-Gdpi=250', '-Grankdir=TB'],
    :medium => ['-Gdpi=250'],
    :large => ['-Gdpi=500'],
  }

  def run(filename, params={})
    size = params[:size] or :large
    system(
      'dot',
      '-Tpng',
      *SIZES[size],
      filename, '-o', output_filename,
    )
  end
end
