module Niconico
  module Helper
    def mylist?(url)
      return true if url =~ /mylist\/\d+/
      false
    end

    def escape_string(str)
      str.gsub(/[\/\\?*:|><]/) { |m| [m.ord + 65_248].pack('U*') }
    end
  end
end
