module Niconico
  module Helper
    def mylist?(url)
      return true if url =~ /mylist\/\d+/
      false
    end
  end
end
