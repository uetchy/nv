require 'thor'

class Nv::CLI < Thor
  include Niconico::Helper

  desc 'dl URL', 'Download video'
  # method_option 'with-comments', :aliases => '-c', :desc => 'Download comments'
  method_option 'with-dir', :aliases => '-d', :desc => 'Create directory'
  method_option 'without-dir', :aliases => '-D', :desc => "Don't create directory"
  def dl(ptr, output=".")
    config = Nv::Config.new(Nv::CONFIG_PATH)
    config.verify_for_authentication!('dl')

    nico = Niconico::Base.new.sign_in(config.email, config.password)

    if mylist?(ptr)
      mylist = nico.mylist(ptr)

      puts "Title : #{mylist.title}"
      puts "Desc  : #{mylist.description}"

      mylist.items.each do |item|
        output = options['without-dir'] ? '.' : escape_string(mylist.title)
        dl(item.link, output)
      end
    else
      video = nico.video(ptr)

      # Inspect
      puts "Downloading... #{video.title}"

      # Donwload video
      output = options['with-dir'] ? escape_string(video.title) : '.'
      video.download output

      # Download comments
      # if options['with-comments']
      #   video.download_comments
      # end

      puts "+ done"
    end
  end

  desc 'info URL', 'Show video/mylist info'
  def info(ptr)
    config = Nv::Config.new(Nv::CONFIG_PATH)
    config.verify_for_authentication!('info')

    nico = Niconico::Base.new.sign_in(config.email, config.password)

    if mylist?(ptr)
      mylist = nico.mylist(ptr)

      puts "Title : #{mylist.title}"
      puts "Desc  : #{mylist.description}"

      mylist.items.each_with_index do |item, i|
        puts "   #{i+1}. #{item.title}"
      end
    else
      video = nico.video(ptr)

      puts video.title
      puts '=' * 40
      puts video.description
      puts '=' * 40
      puts "URL: #{video.watch_url}"
    end
  end

  desc 'browse FILE', 'Open web-browser to show nicovideo page with given movie file'
  def browse(filepath)
    video_id = File.basename(filepath).match(/[^\w]([\w]{2}\d+)[^\w]/)[1]
    system "open http://www.nicovideo.jp/watch/#{video_id}"
  end

  desc 'config KEY VALUE', 'Set config'
  def config(key=nil, value=nil)
    config = Nv::Config.new(Nv::CONFIG_PATH)

    unless key
      puts "=== config(#{Nv::CONFIG_PATH}) ==="
      config.to_h.each do |k, v|
        puts "#{k}=#{v}"
      end
      return
    end

    if value
      config[key] = value
      config.save
    end

    puts "=== config(#{Nv::CONFIG_PATH}) ==="
    puts "#{key}=#{config[key]}"
  end
end
