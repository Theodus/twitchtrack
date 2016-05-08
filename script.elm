module Main (..) where

import Effects exposing (Effects, Never)
import Html exposing (..)
import Html.Shorthand exposing (..)
import Bootstrap.Html exposing (..)
import Http
import Json.Decode as Json exposing (..)
import StartApp
import Task
import Time

type Action
  = NoOp
  | Refresh
  | OnRefresh (Result Http.Error Model)

type alias Model = List Entry

type alias Entry =
  { channel : String
  , game : String
  , stream : String
  , url : String
  , viewers : Int
  }

view : Signal.Address Action -> Model -> Html.Html
view address model =
  container_
    [ h2_ "Twitch Streaming"
    , tableBodyStriped_
      [ thead_
        [ tr_
          [ th_ [ text "Channel" ]
          , th_ [ text "Game" ]
          , th_ [ text "Stream" ]
          ]
        ]
      , tbody_ (List.map rows model)
      ]
    ]

rows : Entry -> Html
rows e =
  tr_
    [ td_ [ a_ e.url e.channel ]
    , td_ [ text (if e.viewers > 0 then e.game else "") ]
    , td_ [ text (if e.viewers > 0 then e.stream else "") ]
    ]

httpTask : Task.Task Http.Error Model
httpTask = Http.get decode "/data"

decode : Json.Decoder Model
decode = at ["channels"] (list decodeEntry)

decodeEntry : Json.Decoder Entry
decodeEntry =
  object5 Entry
    (at [ "channel" ] string)
    (at [ "game" ] string)
    (at [ "stream" ] string)
    (at [ "url" ] string)
    ("viewers" := int)

refreshFx : Effects.Effects Action
refreshFx =
  httpTask
    |> Task.toResult
    |> Task.map OnRefresh
    |> Effects.task

init : (Model, Effects Action)
init = update Refresh []

update : Action -> Model -> (Model, Effects.Effects Action)
update action model =
  case action of
    Refresh ->
      (model, refreshFx)
    OnRefresh result ->
      let
        message = Result.withDefault [] result
      in
        (message, Effects.none)
    _ ->
      (model, Effects.none)

clockSignal : Signal Time.Time
clockSignal = Time.every (Time.minute * 2)

clockRefresh : Signal Action
clockRefresh = Signal.map (\s -> Refresh) clockSignal

app : StartApp.App Model
app =
  StartApp.start
    { init = init
    , inputs = [ clockRefresh ]
    , update = update
    , view = view
    }

main : Signal.Signal Html.Html
main = app.html

port runner : Signal (Task.Task Never ())
port runner = app.tasks
