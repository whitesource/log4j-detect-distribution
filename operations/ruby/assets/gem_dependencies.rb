require 'bundler'
require 'json'

def parse_lock(file_path)
  parser = Bundler::LockfileParser.new(Bundler.read_file(file_path))

  direct = parser.dependencies.keys
  depsToChildren = {}
  deps = {}

  parser.specs.each do |spec|
    children = []
    spec.dependencies.each do |dep|
      children << dep.name
    end
    depsToChildren[spec.name] = children
    deps[spec.name] = {"name": spec.name, "version": spec.version.to_s}
  end

  res = { directDependencies: direct, depsToChildren: depsToChildren, dependencies: deps }
  puts JSON.pretty_generate(res)
end

# ARGV[0] -> path to Gemfile.lock or gems.locked
parse_lock(ARGV[0])