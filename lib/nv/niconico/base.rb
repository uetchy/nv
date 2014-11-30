module Niconico
  class Base < Fabric
    def video(video_id)
      Niconico::Video.new(video_id, @agent)
    end

    def mylist(mylist_id)
      Niconico::Mylist.new(mylist_id, @agent)
    end
  end
end
