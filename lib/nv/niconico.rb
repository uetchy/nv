require 'mechanize'
require 'open-uri'
require 'net/http'
require 'uri'
require 'cgi'
require 'rexml/document'
require 'rss'
require 'ostruct'
require 'ruby-progressbar'

require 'nv/niconico/helper'
require 'nv/niconico/fabric'
require 'nv/niconico/base'
require 'nv/niconico/video'
require 'nv/niconico/mylist'

module Niconico
  OUTPUT_NAME = '%{title} - [%{id}].%{extension}'.freeze
end
