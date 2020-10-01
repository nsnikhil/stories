package test

import (
	"context"
	"fmt"
	newrelic "github.com/newrelic/go-agent"
	"github.com/nsnikhil/stories-proto/proto"
	config2 "github.com/nsnikhil/stories/pkg/config"
	grpc2 "github.com/nsnikhil/stories/pkg/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/alexcesaro/statsd.v2"
	"testing"
	"time"
)

const address = "127.0.0.1:8080"

func TestStories(t *testing.T) {
	go startServer()
	waitForServer()

	//cl := getClient(t)

	//testPingRequest(t, cl)
	//testScenarioOne(t, cl)
}

func testPingRequest(t *testing.T, cl proto.StoriesApiClient) {
	resp, err := cl.Ping(context.Background(), &proto.PingRequest{})
	require.NoError(t, err)

	assert.Equal(t, "pong", resp.GetMessage())
}

func testScenarioOne(t *testing.T, cl proto.StoriesApiClient) {
	ctx := context.Background()

	ar := addRequests()
	sz := len(ar)

	for i := 0; i < sz; i++ {
		resp, err := cl.AddStory(ctx, ar[i])
		require.NoError(t, err)
		assert.True(t, resp.Success)
	}

	mv := mostViewedRequests(sz)

	protoStories := make([]*proto.Story, sz)
	for i, m := range mv {
		resp, err := cl.GetMostViewedStories(ctx, m)
		require.NoError(t, err)

		stories := resp.Stories
		for _, story := range stories {
			protoStories[i] = story
			assert.Equal(t, fmt.Sprintf("Understanding Asymptotic Bounds - Part 1.%d", i+1), story.GetTitle())
		}
	}

	tr := topRatedRequests(sz)
	for i, te := range tr {
		resp, err := cl.GetTopRatedStories(ctx, te)
		require.NoError(t, err)

		stories := resp.Stories
		for _, story := range stories {
			assert.Equal(t, fmt.Sprintf("Understanding Asymptotic Bounds - Part 1.%d", i+1), story.GetTitle())
		}
	}

	gr := getStoryRequests(protoStories)

	for _, g := range gr {
		resp, err := cl.GetStory(ctx, g)
		require.NoError(t, err)

		str := resp.GetStory()

		assert.NotNil(t, str.GetCreatedAtUnix())
		assert.NotNil(t, str.GetUpdatedAtUnix())
		assert.NotNil(t, str.GetTitle())
		assert.NotNil(t, str.GetBody())
	}

	queriesMatrix := [][]string{
		{"analyze", "series", "GO", "IMDb", "AMBIGUOUS", "astronomically", "Logarithmic"},
		{"exponentially", "functions", "Quadratic", "characterize", "improvement", "day"},
		{"consume", "they", "computation", "scale", "easily", "function"},
		{"Binary", "Merge", "Search", "Bubble", "process", "processTwo"},
	}

	time.Sleep(time.Second)

	for _, queries := range queriesMatrix {
		sr := searchRequests(queries)

		for _, r := range sr {
			resp, err := cl.SearchStories(ctx, r)
			require.NoError(t, err)

			stories := resp.Stories

			for _, story := range stories {
				assert.NotNil(t, story.GetCreatedAtUnix())
				assert.NotNil(t, story.GetUpdatedAtUnix())
				assert.NotNil(t, story.GetTitle())
			}
		}
	}

	ur := updateRequests(protoStories)

	for _, r := range ur {
		resp, err := cl.UpdateStory(ctx, r)
		require.NoError(t, err)
		require.True(t, resp.Success)
	}

	gr = getStoryRequests(protoStories)

	for i, g := range gr {
		resp, err := cl.GetStory(ctx, g)
		require.NoError(t, err)

		str := resp.GetStory()

		assert.Equal(t, int64(i+45*2), str.GetViews())
		assert.Equal(t, int64(i+22*4), str.GetUpVotes())
		assert.Equal(t, int64(i+11*4), str.GetDownVotes())
		assert.NotEqual(t, str.GetCreatedAtUnix(), str.GetUpdatedAtUnix())
	}

	dr := deleteRequests(protoStories)

	for _, r := range dr {
		resp, err := cl.DeleteStory(ctx, r)
		require.Nil(t, err)
		require.True(t, resp.GetSuccess())
	}
}

func addRequests() []*proto.AddStoryRequest {
	data := addData()

	sz := len(data)
	addRequests := make([]*proto.AddStoryRequest, sz)

	for i := 0; i < sz; i++ {
		addRequests[i] = &proto.AddStoryRequest{
			Story: &proto.Story{
				Title: data[i].title,
				Body:  data[i].body,
			},
		}
	}

	return addRequests
}

func mostViewedRequests(sz int) []*proto.MostViewedStoriesRequest {
	mostViewedRequests := make([]*proto.MostViewedStoriesRequest, sz)

	for i := 0; i < sz; i++ {
		mostViewedRequests[i] = &proto.MostViewedStoriesRequest{
			Offset: int64(i),
			Limit:  1,
		}
	}

	return mostViewedRequests
}

func topRatedRequests(sz int) []*proto.TopRatedStoriesRequest {
	topRatedRequests := make([]*proto.TopRatedStoriesRequest, sz)

	for i := 0; i < sz; i++ {
		topRatedRequests[i] = &proto.TopRatedStoriesRequest{
			Offset: int64(i),
			Limit:  1,
		}
	}

	return topRatedRequests
}

func getStoryRequests(stories []*proto.Story) []*proto.GetStoryRequest {
	var getStoryRequests []*proto.GetStoryRequest

	for _, story := range stories {
		getStoryRequests = append(getStoryRequests, &proto.GetStoryRequest{
			StoryID: story.GetId(),
		})
	}

	return getStoryRequests
}

func searchRequests(queries []string) []*proto.SearchStoriesRequest {
	var searchRequests []*proto.SearchStoriesRequest

	for _, query := range queries {
		searchRequests = append(searchRequests, &proto.SearchStoriesRequest{
			Query: query,
		})
	}

	return searchRequests
}

func updateRequests(stories []*proto.Story) []*proto.UpdateStoryRequest {
	type updater struct {
		id        string
		views     int
		upVotes   int
		downVotes int
	}

	var updaters []*updater
	for i, story := range stories {
		updaters = append(updaters, &updater{
			id:        story.GetId(),
			views:     i + 45*2,
			upVotes:   i + 22*4,
			downVotes: i + 11*4,
		})
	}

	var updateRequests []*proto.UpdateStoryRequest
	for i, up := range updaters {
		for l := 0; l < up.views; l++ {
			stories[i].Views++
		}

		for l := 0; l < up.upVotes; l++ {
			stories[i].UpVotes++
		}

		for l := 0; l < up.downVotes; l++ {
			stories[i].DownVotes++
		}

		updateRequests = append(updateRequests, &proto.UpdateStoryRequest{
			Story: stories[i],
		})
	}

	return updateRequests
}

func deleteRequests(stories []*proto.Story) []*proto.DeleteStoryRequest {
	var deleteRequests []*proto.DeleteStoryRequest

	for _, story := range stories {
		deleteRequests = append(deleteRequests, &proto.DeleteStoryRequest{
			StoryID: story.GetId(),
		})
	}

	return deleteRequests
}

type testData struct {
	title string
	body  string
}

func addData() []testData {
	return []testData{
		{
			title: "Understanding Asymptotic Bounds - Part 1.1",
			body:  "A postman was asked to deliver atleast 10 letters to their respective location by the evening; He has a bike with a limited amount of fuel. For each delivery, he would earn some fixed amount.\nThe postman also gets a huge bonus if he:\ncan deliver 12 or more letters.\nhas fuel left in the bike.\n\nInstead of picking up 10 random letters and starting the trip, the postman looks at 50 different letters and based on the location, traffic prediction, familiarity, fuel and delivery time chooses 12 letters which maximises his revenue for the day.\nSimilarly, when you are presented with a problem statement, instead of straight away jumping into coding it's always better to first analyze the problem statement, the input data, look for any specific pattern, etc then design the solution.\nOnce you have a solution in hand, how do to verify if your solution is optimal, or how do you figure out is any space for improvement?\nThis is exactly what this and the next few articles in this series are going to cover, we will learn:\nThe tools you can use to analyze your algorithm.\nIn-depth analysis of these tools.\nHow can you use these tools to analyze an algorithm?\nAnalysis of two sorting algorithms from the above learning.\n\nAll the examples in the series are written in GO\n\n\n---\n\nWe will start by looking at the tools we will use to analyze the efficiency of an algorithm in terms of time and space.\nIf you visit websites like Rotten Tomatoes or IMDb you will see ratings for movies, IMDb uses a scale of 1–10 and Rotten Tomatoes uses percentage but at the end, both websites give scores to movies based on various parameters like the story, music, cinematography, etc.\nSimilarly, there exists a scale to rate algorithms, this scale is independent of the hardware, hence it does not have any specific number to rate an algorithm like it would take X-sec or would consume Y-bytes, rather this scale has behaviours; which are used to interpret how an algorithm behaves with respect to a given input, how it affects the time and space taken with respect to a given input.\nWhen we rate an algorithm we bind it to a behaviour in the scale, this bound is an Asymptotic bound.\nWe will now try to understand what the above statement means, but before that let's understand the meaning of word Asymptote.\nWhat does Asmptote mean?\nA straight line that continually approaches a given curve but does not meet it at any finite distance.\nWhy do we care about them?\nWhen we see how a function behaves on a tail end or on an asymptotic end, we have a true understanding of its performance.\nFor example, The code below computes Fibonacci till nth term recursively:",
		},
		{
			title: "Understanding Asymptotic Bounds - Part 1.2",
			body:  "The above function for n = 30 takes around 7 milliseconds and for n = 50 takes 135 seconds, so as the input gets larger and larger the time increases in an exponential manner.\nWhat does Asymptotic bound mean?\nSuppose you design an algorithm to do a certain task, that algorithm can be wrapped up in a function f(x) and you want to understand how a function behaves when the input changes i.e. how the memory or space requirement of the function changes with the respect to the input.\nSo we will pick a function from the scale (scale is described below) to bound it or cover your function and this function is called bounding function and it characterises the behaviour of your function f(x).\n//ADD AN EXAMPLE TO EXPLAIN THE ABOVE AMBIGUOUS PARAGRAPH.\nBounding a given function with other function so as to declare or describe the behaviour of the given function.\nThis bounding function should be as close or as tight as possible to the original function for it to accurately describe the behaviour of the function.\nFor example, given a function f(x) if you can bound the function by another linear function (a function which changes linearly with respect to input) and an exponential function (a function which changes exponentially with respect to input) then the linear function should be the bounding function as it closely resembles the behaviour of the function as you approach the asymptotic end.\nHow do we choose the bounding functions?\nThe bounding scale\nAs mentioned above we pick a value from the scale to interpret the behaviour of our function, this scale has something called as function classes.\nWhat is Function Class?",
		},
		{
			title: "Understanding Asymptotic Bounds - Part 1.3",
			body:  "In the graph above as the value of x increases to an astronomically large number, the tail behaviour or the asymptotic behaviour of the four graphs above remain same i.e they all grow linearly i.e. they have the same behaviour. Hence, we can classify all the function in the graphs as a linear function.\nFor example two functions process and processTwo below does the same exact computation but a bit differently: Though processTwo is doing half the work as process -they both will grow linearly with increase in the input size, i.e. these both functions will have the same asymptotic behaviour, hence they will fall under the same category or same function class.\nSo when we hear a term like linear time/space or factorial time/space or exponential time/space we are not defining a particular function rather we are defining a class of function/behaviours which fall under these umbrellas.\nBelow is the scale which lists well know/used function classes and popular algorithms they bound to:\nConstant Class: Operations like assignment, addition, etc.\nLogarithmic Class: Binary Search, Merge Sort.\nLinear Class: Linear Search.\nQuadratic Class: Bubble Sort.\n\nSo when we choose a bounding function, we tend to choose from a function class, though not necessary your function falls under any of those classes.\n\n\n---\n\nWhen we have to define the behaviour of an algorithm or a function we tend to relate its cost to a know function class like linear or logarithmic, that way we can easily interpret its behaviour when the input gets astronomically large.\n\n\n---\n\nThis wraps the first part of the series, in the next part we will look at the three asymptotic notations and learn how we can use them to characterize our algorithm.",
		},
	}
}

func getClient(t *testing.T) proto.StoriesApiClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	require.NoError(t, err)

	return proto.NewStoriesApiClient(conn)
}

func waitForServer() {
	time.Sleep(time.Second)
}

func startServer() {
	cfg := config2.LoadConfigs()
	lgr := zap.NewNop()

	nrApp, _ := newrelic.NewApplication(newrelic.Config{})

	sc, _ := statsd.New(statsd.Address(cfg.GetStatsDConfig().Address()), statsd.Prefix(cfg.GetStatsDConfig().Namespace()))

	grpc2.NewAppServer(cfg, lgr, nrApp, sc).Start()
}
