module Main exposing (..)

import Models exposing (..)
import Update exposing (..)
import View exposing (..)
import Html exposing (..)


main : Program Never Model Msg
main =
    Html.program
        { init = ( Model "" "", Cmd.none )
        , view = view
        , update = update
        , subscriptions = (\m -> Sub.none)
        }
