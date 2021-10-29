package main

import (
	"fmt"
	"strings"

	"github.com/vikesh-raj/go-textrank/textrank"
)

func main() {

	sentences := []string{
		"Apple Co-Founder Steve Wozniak Launches Company to Clean Up Space Debris.",
		"An Apple co-founder is launching a company to clean up space debris, and not a moment to soon, as NASA estimates there may have been 27,000 pieces of junk floating or hurtling around the Earth last year.",
		"With thousands of additional satellites slated for placement around the Earth as part of future plans for universal internet connectivity, the job to tackle the growing orbital refuse must fall to someone.",
		"The company, called Privateer Space, has nothing to do with piracy, and is in fact in “stealth mode,” and as such we know little about it.",
		"Steve Wozniak, an Apple co-founder of the I software, who has a net worth of $100 million, tweeted out a link to the Privateer website, which currently has nothing but a YouTube video on it.",
		"Later, a press release sent out regarding a titanium 3D printer developed under a company called Desktop Metal featured a quote from Wozniak.",
		`“’3D printing with titanium is incredibly valuable in industries like aerospace because of the material’s ability to support complex and lightweight designs,’ co-founder of Privateer Space, a new satellite company focused on monitoring and cleaning up objects in space,” the press release stated.`,
		"Far from the Musk/Branson/Bezos space race, the fact that at least one tech one-percenter is investing in space clean-up is huge for all of humanity.",
		"Former NASA administrator Jim Bridenstine asked Congress for $15 million solely for a space cleanup mission, tweeting that the ISS had to maneuver out of the way of dangerous space debris on three separate occasions last year.",
		"One of the major problems with space debris is that the smaller it is the more dangerous it becomes, as NASA reports.",
		`“A number of space shuttle windows were replaced because of damage caused by material that was analyzed and shown to be paint flecks,” the agency wrote, noting that debris can travel as fast as 17,500 mph.`,
		`“In fact, millimeter-sized orbital debris represents the highest mission-ending risk to most robotic spacecraft operating in low Earth orbit.”`,
		`Wozniak has said that more information about his cleanup crew will be announced at the AMOS Tech 2021 conference, in Maui, which ended today.`,
	}

	// Summarization
	scores, _ := textrank.SummarizeSentences(sentences, textrank.Options{Debug: true, Language: "english"})
	top := textrank.PickTopSentencesByRatio(scores, 0.3)
	summary := strings.Join(top, "\n")
	fmt.Println("---- Summary ---")
	fmt.Println(summary)

	// Keywords
	text := strings.Join(sentences, " ")
	keywordScores, _ := textrank.Keywords(text, textrank.Options{Language: "english"}, 0.0, 15)
	keywords := textrank.ScoreSentenceToText(keywordScores)
	allKeywords := strings.Join(keywords, "\n")
	fmt.Println("---- Keywords ---")
	fmt.Println(allKeywords)
}
