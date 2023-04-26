package session_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dtimm/roomalone-plugin/pkg/session"
)

var (
	startLocation = session.Location{
		Name:        "start-location",
		Connections: []string{"test-location-1"},
		Description: "start-description",
		Changes:     []string{},
		Story:       "start-story",
	}
	testLocation = session.Location{
		Name:        "test-location-1",
		Connections: []string{"start-location"},
		Description: "test-description-1",
		Changes:     []string{},
		Story:       "test-story-1",
	}
	testLocationUpdated = session.Location{
		Name:        "test-location-1",
		Connections: []string{"start-location"},
		Description: "lo! the description has changed!",
		Changes:     []string{},
		Story:       "test-story-1",
	}
)

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

	Describe("MoveLocation", func() {
		Context("when moving to an existing location", func() {
			It("returns the new location information", func() {
				m, err := testSession.MoveLocation("test-location-1")

				Expect(err).ToNot(HaveOccurred())
				Expect(m).To(BeEquivalentTo(testLocation))

				Expect(testSession.CurrentLocation).To(Equal("test-location-1"))
			})
		})

		Context("with an invalid location", func() {
			It("returns an error", func() {
				_, err := testSession.MoveLocation("fake-location")

				Expect(err).To(MatchError("location fake-location not in map: Move action only accepts location names from the current location connections"))
				Expect(testSession.CurrentLocation).To(Equal("start-location"))
			})
		})
	})

	Describe("GetLocation", func() {
		Context("with empty input", func() {
			It("returns the current location information", func() {
				l, err := testSession.GetLocation("")

				Expect(err).ToNot(HaveOccurred())
				Expect(l).To(BeEquivalentTo(startLocation))
			})
		})

		Context("with an existing location", func() {
			It("returns the location information", func() {
				l, err := testSession.GetLocation("test-location-1")

				Expect(err).ToNot(HaveOccurred())
				Expect(l).To(BeEquivalentTo(testLocation))
			})
		})

		Context("with an invalid location", func() {
			It("returns an error", func() {
				_, err := testSession.GetLocation("fake-location")

				Expect(err).To(MatchError("location fake-location not in map: Location action only accepts existing location names or empty input"))
			})
		})
	})

	Describe("SetLocation", func() {
		Context("with an existing location", func() {
			It("returns the updated location information", func() {
				Expect(testSession.GetLocation("test-location-1")).To(BeEquivalentTo(testLocation))

				l := testSession.SetLocation(testLocationUpdated)

				Expect(l).To(BeEquivalentTo(testLocationUpdated))
				Expect(testSession.GetLocation("test-location-1")).To(BeEquivalentTo(testLocationUpdated))
			})
		})

		Context("with a new location", func() {
			newLocation := session.Location{
				Name:        "new-location",
				Connections: []string{"start-location"},
				Description: "new-description",
				Changes:     []string{},
				Story:       "new-story",
			}
			It("creates the location", func() {
				Expect(testSession.GetLocation("new-location")).Error().To(HaveOccurred())

				l := testSession.SetLocation(newLocation)

				Expect(l).To(BeEquivalentTo(newLocation))
				Expect(testSession.GetLocation("new-location")).To(BeEquivalentTo(newLocation))
			})
		})
	})
})
