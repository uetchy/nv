module Niconico
  class Mylist < Fabric
    include Niconico::Helper

    def initialize(ptr, agent = nil)
      super(agent)

      @id = normalize(ptr)
      @mylist = nil
      fetch
    end

    def method_missing(method, *_args)
      raise(NoMethodError, method) unless @mylist.respond_to? method
      @mylist[method]
    end

    private

    def normalize(ptr)
      ptr.match(/mylist\/([0-9]+)\??/)[1]
    end

    def fetch
      doc = REXML::Document.new(open("http://www.nicovideo.jp/mylist/#{@id}?rss=2.0").read)

      channel = doc.elements['/rss/channel']

      items = []
      channel.elements.each('item') do |item|
        html_description = item.elements['description/text()'].to_s.gsub(/<p class=\"nico-info\">.+<\/p>/, '')
        description = html_description.gsub(/<\/?.*?>/, '')

        items << OpenStruct.new(title: item.elements['title/text()'].to_s,
                                link: item.elements['link/text()'].to_s,
                                guid: item.elements['guid/text()'].to_s,
                                created_at: item.elements['pubDate/text()'].to_s,
                                description: description,
                                html_description: html_description)
      end

      @mylist = OpenStruct.new(title: channel.elements['title/text()'].to_s.gsub(/(^マイリスト\s|‐ニコニコ動画$)/, ''),
                               link: channel.elements['link/text()'].to_s,
                               description: channel.elements['description/text()'].to_s,
                               created_at: channel.elements['pubDate/text()'].to_s,
                               updated_at: channel.elements['lastBuildDate/text()'].to_s,
                               generator: channel.elements['generator/text()'].to_s,
                               author: channel.elements['dc:creator/text()'].to_s,
                               language: channel.elements['language/text()'].to_s,
                               items: items,
                               items_count: items.size)
    end
  end
end
