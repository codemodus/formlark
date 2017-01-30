module Main exposing (..)

import Html exposing (Html, div, input, text)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import String


-- Model


type alias Model =
    { val : String
    }


model : Model
model =
    Model ""



-- Update


type Msg
    = Reverse String


update : Msg -> Model -> Model
update msg model =
    case msg of
        Reverse new ->
            { model | val = new }



-- View


view : Model -> Html Msg
view model =
    div []
        [ input [ placeholder "Text to be reversed.", onInput Reverse ] []
        , div [] [ text (String.reverse model.val) ]
        ]


main : Program Never Model Msg
main =
    Html.beginnerProgram
        { model = model
        , view = view
        , update = update
        }
