module Dashboard.Group exposing
    ( PipelineIndex
    , hdView
    , ordering
    , pipelineDropAreaView
    , pipelineNotSetView
    , view
    )

import Concourse
import Dashboard.DashboardPreview as DashboardPreview
import Dashboard.Drag exposing (drag)
import Dashboard.Group.Models exposing (Group, Pipeline)
import Dashboard.Group.Tag as Tag
import Dashboard.Models exposing (DragState(..), DropState(..))
import Dashboard.Pipeline as Pipeline
import Dashboard.PipelineGridLayout as PipelineGridLayout exposing (cardHeight, cardWidth, padding)
import Dashboard.Styles as Styles
import Dict exposing (Dict)
import HoverState
import Html exposing (Html)
import Html.Attributes exposing (attribute, class, classList, draggable, id, style)
import Html.Events exposing (on, preventDefaultOn, stopPropagationOn)
import Html.Keyed
import Json.Decode
import Maybe.Extra
import Message.Effects as Effects
import Message.Message exposing (DomID(..), Message(..))
import Ordering exposing (Ordering)
import Time
import UserState exposing (UserState(..))
import Views.Spinner as Spinner


ordering : { a | userState : UserState } -> Ordering Group
ordering session =
    Ordering.byFieldWith Tag.ordering (tag session)
        |> Ordering.breakTiesWith (Ordering.byField .teamName)


type alias PipelineIndex =
    Int


type alias Bounds =
    { x : Float, y : Float, width : Float, height : Float }


view :
    { a | userState : UserState }
    ->
        { dragState : DragState
        , dropState : DropState
        , now : Maybe Time.Posix
        , hovered : HoverState.HoverState
        , pipelineRunningKeyframes : String
        , pipelinesWithResourceErrors : Dict ( String, String ) Bool
        , existingJobs : List Concourse.Job
        , pipelines : List Pipeline
        , pipelineLayers : Dict ( String, String ) (List (List Concourse.Job))
        , viewportWidth : Float
        , viewportHeight : Float
        , scrollTop : Float
        }
    -> Group
    -> ( Html Message, Float )
view session params g =
    let
        ( dragTeam, fromIndex, toIndex ) =
            case ( params.dragState, params.dropState ) of
                ( Dragging team fromIdx, NotDropping ) ->
                    ( team, fromIdx, fromIdx + 1 )

                ( Dragging team fromIdx, Dropping toIdx ) ->
                    ( team, fromIdx, toIdx )

                _ ->
                    ( "", -1, -1 )

        pipelinesForGroup =
            if (g.teamName == dragTeam) && (fromIndex >= 0) then
                drag fromIndex toIndex g.pipelines

            else
                g.pipelines

        headerHeight =
            60

        numColumns =
            max 1 (floor (params.viewportWidth / (cardWidth + padding)))

        numRowsVisible =
            max 1 (ceiling ((params.viewportHeight - headerHeight) / (cardHeight + padding)))

        numRowsOffset =
            floor ((params.scrollTop + headerHeight) / (cardHeight + padding))

        isVisible { row, height } =
            (numRowsOffset < row + height)
                && (row <= numRowsOffset + numRowsVisible + 1)

        layersList =
            pipelinesForGroup
                |> List.map
                    (\pipeline ->
                        Dict.get ( pipeline.name, pipeline.teamName ) params.pipelineLayers
                            |> Maybe.withDefault []
                    )

        previewSizes =
            layersList
                |> List.map
                    (\layers ->
                        ( List.length layers
                        , layers
                            |> List.map List.length
                            |> List.maximum
                            |> Maybe.withDefault 0
                        )
                    )

        cards =
            previewSizes
                |> List.map PipelineGridLayout.cardSize
                |> PipelineGridLayout.layout numColumns

        numRows =
            cards
                |> List.map (\c -> c.row + c.height - 1)
                |> List.maximum
                |> Maybe.withDefault 1

        totalCardsHeight =
            numRows
                * cardHeight
                + padding
                * (numRows - 1)

        dropAreaBounds =
            List.range 1 numRows
                |> List.filter (\row -> isVisible { row = row, height = 1 })
                |> List.map toFloat
                |> List.map
                    (\row ->
                        List.range 1 (numColumns + 1)
                            |> List.map toFloat
                            |> List.map
                                (\col ->
                                    { x = (col - 1) * (cardWidth + padding)
                                    , y = (row - 1) * (cardHeight + padding)
                                    , width = 50
                                    , height = cardHeight
                                    }
                                )
                    )
                |> List.concat

        orderings =
            pipelinesForGroup
                |> List.indexedMap (\i _ -> i)

        dropAreaIndexes =
            let
                visibleCards =
                    cards
                        |> List.map2 Tuple.pair orderings
                        |> List.filter (Tuple.second >> isVisible)
            in
            List.range 1 numRows
                |> List.map
                    (\row ->
                        visibleCards
                            |> List.filter (Tuple.second >> .row >> (==) row)
                    )
                |> List.filter (List.isEmpty >> not)
                |> List.map
                    (\rowCards ->
                        let
                            rowHeight =
                                rowCards
                                    |> List.map (Tuple.second >> .height)
                                    |> List.maximum
                                    |> Maybe.withDefault 1

                            maxIndex =
                                rowCards
                                    |> List.map Tuple.first
                                    |> List.maximum
                                    |> Maybe.withDefault 0
                        in
                        rowCards
                            |> List.map (\( i, { width } ) -> List.repeat width i)
                            |> List.concat
                            |> (\i -> i ++ [ maxIndex + 1 ])
                            |> List.repeat rowHeight
                            |> List.concat
                    )
                |> List.concat

        dropAreas =
            dropAreaBounds
                |> List.map2 Tuple.pair dropAreaIndexes
                |> List.map
                    (\( i, bounds ) ->
                        pipelineDropAreaView params.dragState
                            params.dropState
                            g.teamName
                            bounds
                            i
                    )

        pipelineCards =
            if List.isEmpty pipelinesForGroup then
                [ ( "not-set", Pipeline.pipelineNotSetView ) ]

            else
                pipelinesForGroup
                    |> List.map3
                        (\card layers pipeline ->
                            ( card, layers, pipeline )
                        )
                        cards
                        layersList
                    |> List.filter (\( card, _, _ ) -> isVisible card)
                    |> List.map
                        (\( card, layers, pipeline ) ->
                            Html.div
                                [ class "pipeline-wrapper"
                                , style "position" "absolute"
                                , style "transform"
                                    ("translate("
                                        ++ String.fromInt ((card.column - 1) * (cardWidth + padding) + padding)
                                        ++ "px,"
                                        ++ String.fromInt ((card.row - 1) * (cardHeight + padding))
                                        ++ "px)"
                                    )
                                , style
                                    "width"
                                    (String.fromInt
                                        (cardWidth
                                            * card.width
                                            + padding
                                            * (card.width - 1)
                                        )
                                        ++ "px"
                                    )
                                , style "height"
                                    (String.fromInt
                                        (cardHeight
                                            * card.height
                                            + padding
                                            * (card.height - 1)
                                        )
                                        ++ "px"
                                    )
                                ]
                                [ Html.div
                                    ([ class "card"
                                     , style "width" "100%"
                                     , id <| Effects.toHtmlID <| PipelineCard pipeline.id
                                     , attribute "data-pipeline-name" pipeline.name
                                     , attribute
                                        "ondragstart"
                                        "event.dataTransfer.setData('text/plain', '');"
                                     , draggable "true"
                                     , on "dragstart"
                                        (Json.Decode.succeed (DragStart pipeline.teamName pipeline.ordering))
                                     , on "dragend" (Json.Decode.succeed DragEnd)
                                     ]
                                        ++ (if params.dragState == Dragging pipeline.teamName pipeline.ordering then
                                                [ style "width" "0"
                                                , style "margin" "0 12.5px"
                                                , style "overflow" "hidden"
                                                ]

                                            else
                                                []
                                           )
                                        ++ (if params.dropState == DroppingWhileApiRequestInFlight g.teamName then
                                                [ style "opacity" "0.45", style "pointer-events" "none" ]

                                            else
                                                [ style "opacity" "1" ]
                                           )
                                    )
                                    [ Pipeline.pipelineView
                                        { now = params.now
                                        , pipeline = pipeline
                                        , resourceError =
                                            params.pipelinesWithResourceErrors
                                                |> Dict.get ( pipeline.teamName, pipeline.name )
                                                |> Maybe.withDefault False
                                        , existingJobs =
                                            params.existingJobs
                                                |> List.filter
                                                    (\j ->
                                                        j.teamName == pipeline.teamName && j.pipelineName == pipeline.name
                                                    )
                                        , layers = layers
                                        , hovered = params.hovered
                                        , pipelineRunningKeyframes = params.pipelineRunningKeyframes
                                        , userState = session.userState
                                        }
                                    ]
                                ]
                                |> Tuple.pair (String.fromInt pipeline.id)
                        )
    in
    Html.div
        [ id <| Effects.toHtmlID <| DashboardGroup g.teamName
        , class "dashboard-team-group"
        , attribute "data-team-name" g.teamName
        ]
        [ Html.div
            [ style "display" "flex"
            , style "align-items" "center"
            , class <| .sectionHeaderClass Effects.stickyHeaderConfig
            ]
            (Html.div
                [ class "dashboard-team-name" ]
                [ Html.text g.teamName ]
                :: (Maybe.Extra.toList <|
                        Maybe.map (Tag.view False) (tag session g)
                   )
                ++ (if params.dropState == DroppingWhileApiRequestInFlight g.teamName then
                        [ Spinner.spinner { sizePx = 20, margin = "0 0 0 10px" } ]

                    else
                        []
                   )
            )
        , Html.Keyed.node "div"
            [ class <| .sectionBodyClass Effects.stickyHeaderConfig
            , style "position" "relative"
            , style "height" <| String.fromInt totalCardsHeight ++ "px"
            ]
            (pipelineCards ++ [ ( "drop-areas", Html.div [] dropAreas ) ])
        , Html.div
            [ style "position" "absolute"
            , style "top" "0"
            , style "left" "0"
            , style "color" "white"
            , style "width" "100px"
            , style "height" "100px"
            , style "z-index"
                (if g.teamName == "main" then
                    "100000"

                 else
                    "-100000"
                )
            ]
            [ Html.div [] [ Html.text <| "from " ++ String.fromInt fromIndex ]
            , Html.div [] [ Html.text <| "to " ++ String.fromInt toIndex ]
            , Html.div [] [ Html.text <| "indexes " ++ String.fromInt (List.length dropAreaIndexes) ]
            , Html.div [] [ Html.text <| "bounds " ++ String.fromInt (List.length dropAreaBounds) ]
            ]
        ]
        |> (\html ->
                ( html
                , toFloat <| totalCardsHeight + headerHeight
                )
           )


tag : { a | userState : UserState } -> Group -> Maybe Tag.Tag
tag { userState } g =
    case userState of
        UserStateLoggedIn user ->
            Tag.tag user g.teamName

        _ ->
            Nothing


hdView :
    { pipelineRunningKeyframes : String
    , pipelinesWithResourceErrors : Dict ( String, String ) Bool
    , existingJobs : List Concourse.Job
    , pipelines : List Pipeline
    }
    -> { a | userState : UserState }
    -> Group
    -> List (Html Message)
hdView { pipelineRunningKeyframes, pipelinesWithResourceErrors, existingJobs, pipelines } session g =
    let
        pipelinesForGroup =
            pipelines |> List.filter (.teamName >> (==) g.teamName)

        header =
            Html.div
                [ class "dashboard-team-name" ]
                [ Html.text g.teamName ]
                :: (Maybe.Extra.toList <| Maybe.map (Tag.view True) (tag session g))

        teamPipelines =
            if List.isEmpty pipelinesForGroup then
                [ pipelineNotSetView ]

            else
                pipelinesForGroup
                    |> List.map
                        (\p ->
                            Pipeline.hdPipelineView
                                { pipeline = p
                                , pipelineRunningKeyframes = pipelineRunningKeyframes
                                , resourceError =
                                    pipelinesWithResourceErrors
                                        |> Dict.get ( p.teamName, p.name )
                                        |> Maybe.withDefault False
                                , existingJobs = existingJobs
                                }
                        )
    in
    case teamPipelines of
        [] ->
            header

        p :: ps ->
            -- Wrap the team name and the first pipeline together so
            -- the team name is not the last element in a column
            Html.div
                (class "dashboard-team-name-wrapper" :: Styles.teamNameHd)
                (header ++ [ p ])
                :: ps


pipelineNotSetView : Html Message
pipelineNotSetView =
    Html.div
        [ class "card" ]
        [ Html.div
            Styles.noPipelineCardHd
            [ Html.div
                Styles.noPipelineCardTextHd
                [ Html.text "no pipelines set" ]
            ]
        ]


pipelineDropAreaView : DragState -> DropState -> String -> Bounds -> Int -> Html Message
pipelineDropAreaView dragState dropState name { x, y, width, height } index =
    let
        active =
            case dragState of
                Dragging team _ ->
                    team == name

                _ ->
                    False
    in
    Html.div
        [ classList
            [ ( "drop-area", True )
            , ( "active", active )
            , ( "animation", dropState /= NotDropping )
            ]
        , style "position" "absolute"
        , style "transform" <|
            "translate("
                ++ String.fromFloat x
                ++ "px,"
                ++ String.fromFloat y
                ++ "px)"
        , style "width" <| String.fromFloat width ++ "px"
        , style "height" <| String.fromFloat height ++ "px"
        , on "dragenter" (Json.Decode.succeed (DragOver name index))

        -- preventDefault is required so that the card will not appear to
        -- "float" or "snap" back to its original position when dropped.
        , preventDefaultOn "dragover" (Json.Decode.succeed ( DragOver name index, True ))
        , stopPropagationOn "drop" (Json.Decode.succeed ( DragEnd, True ))
        , attribute "data-index" (String.fromInt index)
        ]
        [ Html.text (String.fromInt index) ]
