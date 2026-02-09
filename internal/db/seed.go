package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/high-la/gopher-social/internal/store"
)

var usernames = []string{

	"Alice", "Blake", "Caleb", "Diana", "Ethan", "Fiona", "George",
	"Hannah", "Isaac", "James", "Kevin", "Laura", "Mason", "Nora",
	"Oliver", "Paula", "Quinn", "Riley", "Sarah", "Thomas", "Uriel",
	"Victor", "Wyatt", "Xander", "Yara", "Zane",
	"Aaron", "Beatrice", "Cody", "Dakota", "Elena", "Felix", "Gideon",
	"Hazel", "Iris", "Jasper", "Kira", "Leo", "Maya", "Nolan", "Owen",
	"Phoebe", "Quentin", "Ruby", "Silas", "Tessa", "Uma", "Vera",
	"Willa", "Xavier", "Yosef",
}

var titles = []string{
	"Hidden Gems", "Morning Routine", "Fast Growth", "Digital Detox",
	"Daily Habits", "Simple Living", "Tech Trends", "Deep Focus",
	"Travel Hacks", "Mindset Shift", "Future Goals", "Strict Rules",
	"Modern Art", "Quiet Power", "Smart Saving", "Urban Life",
	"Hard Truths", "Fresh Starts", "Core Values", "Inner Peace",
}

var contents = []string{
	"Technology moves fast. To stay relevant, you must keep learning every single day. Adapt or get left behind in this digital age.",
	"Constant distractions kill your productivity. Set aside blocks of time for focused, uninterrupted work. That is where the real progress happens.",
	"Writing clever code is easy, but writing readable code is hard. Aim for simplicity so others can maintain your work. Your future self will thank you.",
	"Growth requires stepping out of your comfort zone. Embrace the friction of new challenges and learn from every mistake you make along the way.",
	"Consistency is the secret sauce to success. Small, repetitive actions build into massive results over several months. Just show up every single morning.",
	"A solid foundation makes scaling a breeze. Separate your concerns and keep your logic decoupled from your database. This ensures your app stays flexible.",
	"Automation is changing the job market rapidly. Focus on high-level problem solving and emotional intelligence. These are the skills machines cannot easily replace.",
	"Your attention is your most valuable currency. Protect it by turning off notifications and clearing your workspace. One task at a time leads to mastery.",
	"Less is always more in user experience. Remove the clutter until only the essential features remain. A clean interface creates a happy, returning user.",
	"The best code often goes unnoticed by the user. It handles edge cases and errors silently in the background. Stability is the hallmark of a great engineer.",
	"Building for millions requires a different mindset. Think about caching, load balancing, and database sharding early on. Prepare your infrastructure for the heavy traffic.",
	"Know exactly why you are building a product. Clear goals lead to better decision making and faster execution. Never lose sight of the primary user problem.",
	"Do not over-engineer the simple parts of your app. Keep the business logic central and easy to find in your files. Clarity beats complexity every single time.",
	"Remote work has changed how we collaborate forever. Use asynchronous communication to keep the team moving without constant meetings. Trust your developers to deliver.",
	"Your database is the heart of your application. Choose your schema wisely and index your columns for speed. A slow base layer ruins the entire experience.",
	"Entering the flow state is a developer's superpower. Minimize interruptions to stay in the zone for hours. This is when the most complex bugs get solved.",
	"Starting a new project is the hardest part. Don't wait for the perfect idea to strike you. Just write the first line of code and keep going.",
	"Production environments demand extreme caution and care. Always use staging environments to test your latest changes. One wrong command can take down the whole system.",
	"Think about how data flows through your entire stack. Map out the connections between services before you start coding. A good plan prevents a thousand future headaches.",
	"The industry changes every few months with new tools. Don't be afraid to pivot your strategy when a better way emerges. Flexibility is a survival trait for developers.",
}

var tags = []string{
	"Go", "Backend", "Programming", "Tech", "Database", "Software",
	"WebDev", "Cloud", "API", "Scalability", "Architecture", "Coding",
	"PostgreSQL", "Systems", "Efficiency", "Logic", "DevOps",
	"Microservices", "Design", "Performance",
}

var comments = []string{
	"Great read, thanks for sharing!", "I totally agree with this point.",
	"Could you explain this further?", "This helped me a lot today.",
	"Very well written and clear.", "I have a different perspective.",
	"Thanks for the useful tips!", "Interesting take on the topic.",
	"Adding this to my bookmarks.", "Saved me hours of research.",
	"Short, sweet, and to the point.", "The architecture part is spot on.",
	"Would love to see a part two.", "Simple and very effective.",
	"Nice work on the breakdown.", "I'm definitely trying this out.",
	"Exactly what I was looking for.", "Keep up the great content!",
	"Mind-blowing perspective here.", "Concise and very professional.",
}

func Seed(store store.Storage) {

	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {

		err := store.Users.Create(ctx, user)
		if err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	// ,
	posts := generatePosts(200, users)
	for _, post := range posts {

		err := store.Posts.Create(ctx, post)
		if err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	// ,
	comments := generateComments(500, users, posts)
	for _, comment := range comments {

		err := store.Comments.Create(ctx, comment)
		if err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("seeding completed")
}

func generateUsers(num int) []*store.User {

	users := make([]*store.User, num)

	for i := range num {

		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "1234",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {

	posts := make([]*store.Post, num)

	for i := range num {

		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {

	cms := make([]*store.Comment, num)

	for i := range num {

		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms

}
