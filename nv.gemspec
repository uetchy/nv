# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'nv/version'

Gem::Specification.new do |spec|
  spec.name          = "nv"
  spec.version       = Nv::VERSION
  spec.authors       = ["Yasuaki Uechi"]
  spec.email         = ["uetchy@randompaper.co"]
  spec.summary       = %q{The toolbelt for nicovideo}
  spec.description   = %q{The commandline tool for downloading videos and mylist at nicovideo.}
  spec.homepage      = ""
  spec.license       = "MIT"

  spec.files         = `git ls-files -z`.split("\x0")
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.7"
  spec.add_development_dependency "rake", "~> 10.0"

  spec.add_dependency "activesupport"
  spec.add_dependency "mechanize"
  spec.add_dependency "thor"
  spec.add_dependency "ruby-progressbar"
end
