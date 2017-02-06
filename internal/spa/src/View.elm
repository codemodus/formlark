module View exposing (view)

import Update exposing (..)
import Models exposing (Model)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)


view : Model -> Html Msg
view model =
    let
        theLogin : String -> Html Msg -> Html Msg
        theLogin title content =
            div [ class "column" ]
                [ div [ class "card" ]
                    [ header [ class "card-header" ]
                        [ p [ class "card-header-title" ]
                            [ text title ]
                        , span [ class "card-header-icon" ]
                            [ span [ class "icon" ]
                                [ i [ class "fa fa-file-text" ]
                                    []
                                ]
                            ]
                        ]
                    , div [ class "card-content" ]
                        [ div [ class "content" ]
                            [ content ]
                        ]
                    ]
                ]

        regForm =
            Html.form [ action "/user", method "POST", id "form" ]
                [ label [ class "label" ] [ text "Email" ]
                , email
                , label [ class "label" ] [ text "Domain Name" ]
                , domain
                , input [ type_ "submit", value "Get my form!", class "button is-primary" ] []
                ]

        example =
            p []
                [ pre [ style [ ( "padding", ".6em" ) ] ]
                    [ text "<form method=\"POST\" action=\"https://formlark.com/key/<key>/send\"> \n\t<input type=\"text\" name=\"message\" /> \n\t<input type=\"submit\" /> \n</form>"
                    ]
                ]

        email =
            p [ class "control has-icon has-icon-right is-large" ]
                [ input [ class "input is-medium", placeholder "you@example.com", type_ "email", onInput UpdateEmail ]
                    []
                ]

        domain =
            p [ class "control" ]
                [ input [ class "input is-medium", placeholder "example.com", type_ "text", onInput UpdateDomain ]
                    []
                ]

        hero : String -> String -> Html Msg
        hero title sub =
            section [ class "hero" ]
                [ div [ class "hero-body" ]
                    [ div [ class "container" ]
                        [ h1 [ class "title" ]
                            [ text title ]
                        , h2 [ class "subtitle" ]
                            [ text sub ]
                        ]
                    ]
                ]

        hero_ : String -> Html Msg -> Html Msg
        hero_ title body =
            section [ class "hero is-primary" ]
                [ div [ class "hero-body" ]
                    [ div [ class "container" ]
                        [ h1 [ class "title" ]
                            [ text title ]
                        , body
                        ]
                    ]
                ]
    in
        div []
            [ div [ class "container section" ]
                [ div [ class "columns" ]
                    [ div [ class "column" ]
                        [ header []
                            [ hero "Form Lark" "Dynamic forms for static sites without exposing your email."
                            ]
                        , div [ class "columns" ]
                            [ div [ class "column" ] []
                            , theLogin "Get a form" regForm
                            , div [ class "column" ] []
                            ]
                        ]
                    ]
                ]
            , hero_ "Example" example
            ]
