module Update exposing (..)

import Models exposing (Model)


type Msg
    = UpdateEmail String
    | UpdateDomain String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdateEmail newEmail ->
            ( { model | email = newEmail }, Cmd.none )

        UpdateDomain newDomain ->
            ( { model | domain = newDomain }, Cmd.none )
