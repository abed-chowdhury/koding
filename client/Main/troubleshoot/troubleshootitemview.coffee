class TroubleshootItemView extends KDCustomHTMLView

  constructor: (options, data) ->
    options.cssClass ?= "hidden troubleshootitem"
    super options, data

    @loader = new KDLoaderView
      size          : width : 16
      showLoader    : yes

    @getData().once "healthCheckStarted", @bound "startCheck"

    @getData().once "healthCheckCompleted", @bound "completeCheck"

    @forwardEvent @getData(), "recoveryStarted"
    @on "recoveryStarted", @bound "startCheck"

    @forwardEvent @getData(), "recoveryCompleted"
    @on "recoveryCompleted", @bound "completeCheck"


    @status = new KDCustomHTMLView
      tagName   : "strong"
      pistachio : "{{#(status)}}"
      cssClass  : "status"
    , @getData()

    {title} = @getOptions()
    @title = new KDCustomHTMLView
      tagName  : "span"
      partial  : "#{title}"
      cssClass : "checker-name"

  viewAppended: JView::viewAppended

  # instead we need a slowness indicator
  getResponseTime: ->
    responseTime = @getData().getResponseTime()
    unless responseTime is 0
      return "#{responseTime} ms"
    else
      ""

  pistachio:->
    """
    {{> @title}}{{> @status }}{{> @loader}}
    """

  startCheck: ->
    @show()
    @loader.show()
    @status.render()
    @status.unsetClass "fail success"

  completeCheck: ->
    @loader.hide()
    @status.render()
    {status} = @getData()
    @status.unsetClass "fail success"
    @status.setClass "#{status}"
