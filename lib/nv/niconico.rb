require 'mechanize'
require 'open-uri'
require 'net/http'
require 'uri'
require 'cgi'
require 'rexml/document'
require 'rss'
require 'ostruct'
require 'ruby-progressbar'

$LOAD_PATH.unshift(File.expand_path(File.dirname(__FILE__)))

require 'niconico/helper'
require 'niconico/fabric'
require 'niconico/base'
require 'niconico/video'
require 'niconico/mylist'

module Niconico
  OUTPUT_NAME = "%{title} - [%{id}].%{extension}"
end
