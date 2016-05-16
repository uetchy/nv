require 'ostruct'
require 'yaml'
require 'fileutils'

class String
  def undent
    gsub(/^.{#{slice(/^ +/).length}}/, '')
  end
end

module Nv
  class Config < OpenStruct
    def initialize(config_path)
      # Initialize config file
      @config_path = config_path
      config_dir = File.dirname(@config_path)

      Dir.mkdir(config_dir) unless Dir.exist?(config_dir)
      FileUtils.touch(@config_path) unless File.exist?(@config_path)

      @config = YAML.load(open(@config_path).read) || {}
      super(@config)
    end

    def save
      File.open(@config_path, 'w') do |f|
        f.print YAML.dump(transform_keys(to_h, &:to_s))
      end
    end

    def verify_for_authentication
      email && password
    end

    def verify_for_authentication!(cmd)
      unless verify_for_authentication
        puts <<-EOD.undent
        `nv #{cmd}` should be given email and password.
        $ nv config email <email>
        $ nv config password <password>
        EOD
        exit
      end
    end

    private

    def transform_keys(hs)
      result = {}
      hs.each_key do |key|
        result[yield(key)] = hs[key]
      end
      result
    end
  end
end
