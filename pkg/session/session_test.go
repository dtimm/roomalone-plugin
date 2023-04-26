package session_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dtimm/roomalone-plugin/pkg/session"
)

// const (
// 	startLocation = `{
//         "name": "start-location",
//         "connections": ["test-location-2"],
//         "description": "start-description",
//         "changes": [],
//         "story": "start-story"
//     }`
// 	testLocation = `{
// 		"name": "test-location-1",
// 		"connections": ["start-location"],
// 		"description": "test-description-1",
// 		"changes": [],
// 		"story": "test-story-1"
// 	}`
// 	testLocationWithChanges = `{
// 		"name": "test-location-1",
// 		"connections": ["start-location"],
// 		"description": "test-description-1",
// 		"changes": ["test-change-1"],
// 		"story": "test-story-1"
// 	}`
// 	testInventory = `{"items":["test-item-1","test-item-2"]}`
// )

var _ = Describe("Session", func() {
	var testSession *session.Session

	BeforeEach(func() {
		testSession = session.New("../../story/test_game")
	})

	Describe("New", func() {
		Context("with a valid game name", func() {
			It("should load the files for the game", func() {
				Expect(testSession.CurrentLocation).To(Equal("start-location"))
				Expect(testSession.Items).To(HaveLen(2))
				Expect(testSession.Locations).To(HaveLen(2))
			})
		})
	})

	Describe("SetLocation", func() {
		// Context("when moving to an existing location", func() {
		// 	It("returns the new location information", func() {
		// 		m, err := testSession.SetLocation("test-location-1")

		// 		Expect(err).ToNot(HaveOccurred())
		// 		Expect(m.Role).To(Equal(openai.ChatMessageRoleUser))
		// 		Expect(m.Content).To(MatchJSON(testLocation))
		// 	})
		// })

		Context("with an invalid location", func() {
			It("returns an error", func() {
				_, err := testSession.SetLocation("fake-location")

				Expect(testSession.CurrentLocation).To(Equal("start-location"))
				Expect(err).To(MatchError("location fake-location not in map: Move action only accepts location names from the current location connections"))
			})
		})
	})

	Describe("GetLocation", func() {
		// Context("with empty input", func() {
		// 	It("returns the current location information", func() {
		// 		l, err := testSession.GetLocation("")

		// 		Expect(err).ToNot(HaveOccurred())
		// 		Expect(l.Role).To(Equal(openai.ChatMessageRoleUser))
		// 		Expect(l.Content).To(MatchJSON(startLocation))
		// 	})
		// })

		// Context("with an existing location", func() {
		// 	It("returns the location information", func() {
		// 		l, err := testSession.GetLocation("test-location-1")

		// 		Expect(err).ToNot(HaveOccurred())
		// 		Expect(l.Role).To(Equal(openai.ChatMessageRoleUser))
		// 		Expect(l.Content).To(MatchJSON(testLocation))
		// 	})
		// })

		Context("with an invalid location", func() {
			It("returns an error", func() {
				_, err := testSession.GetLocation("fake-location")

				Expect(err).To(MatchError("location fake-location not in map: Location action only accepts existing location names or empty input"))
			})
		})
	})

	Describe("AddLocationChange", func() {
		// Context("with an existing location", func() {
		// 	It("returns the updated location information", func() {
		// 		l, err := testSession.AddLocationChange("test-location-1: test-change-1")

		// 		Expect(err).ToNot(HaveOccurred())
		// 		Expect(l.Role).To(Equal(openai.ChatMessageRoleUser))
		// 		Expect(l.Content).To(MatchJSON(testLocationWithChanges))
		// 	})
		// })

		Context("with an invalid location", func() {
			It("returns an error", func() {
				_, err := testSession.AddLocationChange("fake-location: test-change-1")

				Expect(err).To(MatchError("location fake-location not in map: only use location names that you know exist"))
			})
		})

		Context("with incorrectly formatted input", func() {
			It("returns an error", func() {
				_, err := testSession.AddLocationChange("add a change to fake-location")

				Expect(err).To(MatchError("input `add a change to fake-location` did not match format: `location: description of change`"))
			})
		})
	})
})
