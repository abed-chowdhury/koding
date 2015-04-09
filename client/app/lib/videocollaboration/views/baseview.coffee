kd = require 'kd'
getNick = require 'app/util/nick'
VideoControlView = require './controlview'
isMyChannel = require 'app/util/isMyChannel'

module.exports = class ChatVideoView extends kd.View

  constructor: (options = {}, data) ->

    options.cssClass = kd.utils.curry 'ChatVideo', options.cssClass

    super options, data

    @publisherView = null

    @createContainer()
    @createOfflineUserContainer()
    @createControls()


  createContainer: ->

    @container = new kd.CustomHTMLView { cssClass: 'ChatVideo-container' }
    @addSubView @container


  createOfflineUserContainer: ->

    @offlineUserContainer = new kd.CustomHTMLView { cssClass: 'ChatVideo-offlineUserContainer' }
    @addSubView @offlineUserContainer


  createControls: ->

    @controls = new kd.CustomHTMLView { cssClass: 'ChatVideo-controls' }

    @controlAudio = createVideoControl 'audio', no
    @controlAudio.on 'ActiveStateChangeRequested', @handleStateChangeRequest 'audio'
    @controls.addSubView @controlAudio

    @controlVideo = createVideoControl 'video', no
    @controlVideo.on 'ActiveStateChangeRequested', @handleStateChangeRequest 'video'
    @controls.addSubView @controlVideo

    if isMyChannel @getData()
      @controlEnd = createVideoControl 'end', no
      @controlEnd.on 'ActiveStateChangeRequested', @handleStateChangeRequest 'end'
      @controls.addSubView @controlEnd

    @addSubView @controls


  getContainer: -> @container


  getOfflineUserContainer: -> @offlineUserContainer


  ###*
   * This method needs to be overriden by subclasses.
   *
   * @abstract
  ###
  handleStateChangeRequest: (type) -> (active) -> throw new Error 'needs to be implemented'


  show: ->

    @emit 'ViewWillShow'

    super

    @emit 'ViewDidShow'


  hide: ->

    @emit 'ViewWillHide'

    super

    @emit 'ViewDidHide'


createVideoControl = (type, active) ->

  cssClass = "ChatVideoControl ChatVideoControl--#{type}"

  title = switch type
    when 'video'
      activeTooltipText = 'Turn-off Camera'
      deactiveTooltipText = 'Turn-on Camera'
    when 'audio'
      activeTooltipText = 'Mute'
      deactiveTooltipText = 'Unmute'
    when 'end'
      activeTooltipText = 'End session'
      deactiveTooltipText = 'End session'

  new VideoControlView { cssClass, active, activeTooltipText, deactiveTooltipText }


