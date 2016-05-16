module Niconico
  class Video < Fabric
    include Niconico::Helper

    def initialize(ptr, agent = nil)
      super(agent)

      @thumb = nil
      @flv   = nil
      @id    = normalize(ptr) # pvid | thumb_cached
    end

    def download(output = '.')
      escaped_title  = escape_string(thumb.title)
      escaped_output = escape_string(output)

      filename = sprintf OUTPUT_NAME, title: escaped_title, id: thumb.video_id, extension: thumb.extension
      filepath = File.join(escaped_output, filename)
      filepath_nvdownload = "#{filepath}.nvdownload"

      # Return when video file is already exist
      if File.exist? filepath
        File.delete filepath_nvdownload if File.exist? filepath_nvdownload
        return
      end

      # Create output dir
      Dir.mkdir escaped_output unless Dir.exist? escaped_output

      # Define request headers
      options = {
        'Cookie' => flv.history_cookies
      }
      if File.exist? filepath_nvdownload
        options['Range'] = "bytes=#{File.size(filepath_nvdownload)}-"
      end

      # Download video
      progress_bar = nil
      url = URI.parse(flv.url)
      begin
        Net::HTTP.start(url.host, url.port) do |http|
          header = http.request_head("#{url.path}?#{url.query}", options)
          progress_bar = ProgressBar.create(total: header['content-length'].to_i)
          transferred_bytes = 0
          request = Net::HTTP::Get.new(url, options)
          http.request request do |response|
            open(filepath_nvdownload, 'ab') do |io|
              response.read_body do |chunk|
                transferred_bytes += chunk.size
                if progress_bar
                  progress_bar.progress = transferred_bytes
                else
                  puts "#{transferred_bytes} / Total size is unknown"
                end

                io.write chunk
              end
            end
          end
        end
      rescue => e
        puts "Failed download: #{e}"
        return
      end

      # Rename .nvdownload to real file
      File.rename(filepath_nvdownload, filepath)
    end

    # GET http://flapi.nicovideo.jp/api/getwaybackkey?thread=1345476375
    # => waybackkey=1417346808.E9d0LUF9gvFvt3Rrf5TP91Pa0LA

    # POST http://msg.nicovideo.jp/53/api/
    # <packet><thread thread="1345476375" version="20090904" user_id="1501297" scores="1" nicoru="1" with_global="1"/><thread_leaves thread="1345476375" user_id="1501297" scores="1" nicoru="1">0-14:100,1000</thread_leaves></packet>
    # <packet>
    #   <thread thread="1345476375"
    #           version="20090904"
    #           waybackkey="1417346808.E9d0LUF9gvFvt3Rrf5TP91Pa0LA"
    #           when="1417346804"
    #           user_id="1501297"
    #           scores="1"
    #           nicoru="1"
    #   />
    #   <thread_leaves
    #           thread="1345476375"
    #           waybackkey="1417346808.E9d0LUF9gvFvt3Rrf5TP91Pa0LA"
    #           when="1417346804" # 起点
    #           user_id="1501297"
    #           res_before="8541"
    #           scores="1"
    #           nicoru="1" >0-14:100,1000</thread_leaves>
    # </packet>
    # <chat thread="1345476375" no="8540" vpos="55658" date="1417346426" mail="184" user_id="tzaiW5hp-SvJG6UGkEa0kELQd3w" anonymity="1" leaf="9">白いレースのハンカチかな？</chat>
    # <chat thread="1345476375" no="8539" vpos="41534" date="1417233687" mail="184" user_id="1p_4U9fr-YvliRaUl3E6XGEavp4" premium="1" anonymity="1" leaf="6">くっそｗｗｗｗ</chat>
    def download_comments(output = '.')
      escaped_title  = escape_string(thumb.title)
      escaped_output = escape_string(output)

      filename = sprintf OUTPUT_NAME, title: escaped_title, id: thumb.video_id, extension: 'comments'
      filepath = File.join(escaped_output, filename)

      Dir.mkdir(escaped_output) unless Dir.exist? escaped_output

      url = URI.parse(flv.ms)

      thread_id = flv.thread_id
      length = (flv.l.to_i / 60).round
      res = Net::HTTP.new(url.host).start do |http|
        req = Net::HTTP::Post.new(url.path, 'Cookie' => flv.history_cookies)
        req.body = %(<packet><thread thread="#{thread_id}" version="20090904" scores="1" nicoru="1" with_global="1"/><thread_leaves thread="#{thread_id}" scores="1" nicoru="1">0-#{length}:10</thread_leaves></packet>)
        http.request(req)
      end

      open(filepath, 'w') do |f|
        f.write res.body
      end

      # TODO: Comment parser
      # doc = REXML::Document.new(res.body)
      # chats = doc.elements.to_a('/packet/chat')
      # chats.each do |chat|
      #   puts chat.attribute('vpos')
      #   puts chat.text
      # end
    end

    # Combined parameter fetcher
    def method_missing(method, *_args)
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

    attr_reader :id

    def thumb(id = @id)
      return @thumb if @thumb

      doc = REXML::Document.new(open("http://ext.nicovideo.jp/api/getthumbinfo/#{id}"))
      watch_url = doc.elements['nicovideo_thumb_response/thumb/watch_url'].text.to_s
      perm_video_id = strip_id(watch_url)

      @thumb = OpenStruct.new(video_id: doc.elements['nicovideo_thumb_response/thumb/video_id'].text.to_s,
                              perm_video_id: perm_video_id,
                              watch_url: watch_url,
                              extension: doc.elements['nicovideo_thumb_response/thumb/movie_type'].text.to_s,
                              title: doc.elements['nicovideo_thumb_response/thumb/title'].text.to_s,
                              description: doc.elements['nicovideo_thumb_response/thumb/description'].text.to_s)
    end

    def flv(id = @id)
      return @flv if @flv
      @agent.get("http://www.nicovideo.jp/watch/#{id}")
      history_cookie = @agent.cookies.map(&:to_s).join('; ')

      flv_str = @agent.get("http://www.nicovideo.jp/api/getflv?v=#{id}").body
      flv_hash = Hash[flv_str.split('&').map { |e| i = e.split('='); i[1] = CGI.unescape(i[1]); i }]
      flv_hash.update(history_cookies: history_cookie)
      @flv = OpenStruct.new(flv_hash)
    end
  end
end
