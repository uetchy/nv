module Niconico
  class Video < Fabric
    include Niconico::Helper

    def initialize(ptr, agent=nil)
      super(agent)

      @thumb = nil
      @flv   = nil
      @id    = normalize(ptr) # pvid | thumb_cached
    end

    def download(output=".")
      escapedTitle = thumb.title.gsub(/\//, "／")
      filename = "#{escapedTitle} - [#{thumb.video_id}].#{thumb.extension}"
      filepath = File.join(output, filename)

      Dir.mkdir(output) unless Dir.exist? output

      progress_bar = nil
      File.open(filepath, 'wb') do |fp|
        open(flv.url, 'rb',
          'Cookie' => flv.history_cookies,
          :content_length_proc => lambda{ |content_length|
            if content_length
              progress_bar = ProgressBar.create(:total => content_length)
            end
          },
          :progress_proc => lambda{ |transferred_bytes|
            if progress_bar
              progress_bar.progress = transferred_bytes
            else
              puts "#{transferred_bytes} / Total size is unknown"
            end
          }
        ) do |f|
          fp.print f.read
        end
      end
    end

    # Combined parameter fetcher
    def method_missing(method, *args)
      thumb[method] || flv[method] || raise(NoMethodError, method)
    end

    private

    def strip_id(url)
      url.match(/(?:watch\/)?(\w{2}?\d+)/)[1]
    end

    def normalize(ptr)
      vid = strip_id(ptr)
      thumb(vid).perm_video_id unless vid =~ /\A\d+\z/
    end

    def id
      @id
    end

    def thumb(id=@id)
      return @thumb if @thumb

      doc = REXML::Document.new(open("http://ext.nicovideo.jp/api/getthumbinfo/#{id}"))
      watch_url = doc.elements['nicovideo_thumb_response/thumb/watch_url'].text.to_s
      perm_video_id = strip_id(watch_url)

      @thumb = OpenStruct.new({
        :video_id    => doc.elements['nicovideo_thumb_response/thumb/video_id'].text.to_s,
        :perm_video_id => perm_video_id,
        :watch_url   => watch_url,
        :extension   => doc.elements['nicovideo_thumb_response/thumb/movie_type'].text.to_s,
        :title       => doc.elements['nicovideo_thumb_response/thumb/title'].text.to_s,
        :description => doc.elements['nicovideo_thumb_response/thumb/description'].text.to_s
      })
    end

    def flv(id=@id)
      return @flv if @flv

      @agent.get("http://www.nicovideo.jp/watch/#{id}")
      history_cookie = @agent.cookies.map(&:to_s).join('; ')

      flv_str = @agent.get("http://www.nicovideo.jp/api/getflv?v=#{id}").body
      flv_hash = Hash[flv_str.split("&").map{|e| i=e.split('=');i[1]=CGI.unescape(i[1]);i }]
      flv_hash.update(:history_cookies => history_cookie)
      @flv = OpenStruct.new(flv_hash)
    end
  end
end
