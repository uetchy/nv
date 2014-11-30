require "nv/version"
require "nv/config"
require "nv/niconico"
require "nv/cli"

=begin

nico = Niconico.new.sign_in(...)
video = nico.video('sm9')
video.download('./')

mylist = nico.mylist('482029')
mylist.download()

include Niconico::Helper
if mylist? 'http://...'
  ...
end

   ##### Way #####

1. nico = Niconico::Base.new.sign_in(...)
   video = nico.video('sm9') => Niconico::Video
   puts video.id
   video.download

2. video = Niconico::Video.new('sm9')
   video.sign_in(...)
   video.download

=end

module Nv
  class LackOfInformation < StandardError; end

  CONFIG_PATH = File.join(ENV['HOME'], '.config', 'nv')
end
