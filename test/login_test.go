package test

import (
	"time"

	"github.com/onsi/ginkgo/example/books"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checking books out of the library", Label("library"), func() {
	var library *libraries.Library
	var book *books.Book
	var valjean *users.User
	BeforeEach(func() {
		library = libraries.NewClient()
		book = &books.Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
		}
		valjean = users.NewUser("Jean Valjean")
	})

	When("the library has the book in question", func() {
		BeforeEach(func(ctx SpecContext) {
			Expect(library.Store(ctx, book)).To(Succeed())
		})

		Context("and the book is available", func() {
			It("lends it to the reader", func(ctx SpecContext) {
				Expect(valjean.Checkout(ctx, library, "Les Miserables")).To(Succeed())
				Expect(valjean.Books()).To(ContainElement(book))
				Expect(library.UserWithBook(ctx, book)).To(Equal(valjean))
			}, SpecTimeout(time.Second*5))
		})

		Context("but the book has already been checked out", func() {
			var javert *users.User
			BeforeEach(func(ctx SpecContext) {
				javert = users.NewUser("Javert")
				Expect(javert.Checkout(ctx, library, "Les Miserables")).To(Succeed())
			})

			It("tells the user", func(ctx SpecContext) {
				err := valjean.Checkout(ctx, library, "Les Miserables")
				Expect(err).To(MatchError("Les Miserables is currently checked out"))
			}, SpecTimeout(time.Second*5))

			It("lets the user place a hold and get notified later", func(ctx SpecContext) {
				Expect(valjean.Hold(ctx, library, "Les Miserables")).To(Succeed())
				Expect(valjean.Holds(ctx)).To(ContainElement(book))

				By("when Javert returns the book")
				Expect(javert.Return(ctx, library, book)).To(Succeed())

				By("it eventually informs Valjean")
				notification := "Les Miserables is ready for pick up"
				Eventually(ctx, valjean.Notifications).Should(ContainElement(notification))

				Expect(valjean.Checkout(ctx, library, "Les Miserables")).To(Succeed())
				Expect(valjean.Books(ctx)).To(ContainElement(book))
				Expect(valjean.Holds(ctx)).To(BeEmpty())
			}, SpecTimeout(time.Second*10))
		})
	})

	When("the library does not have the book in question", func() {
		It("tells the reader the book is unavailable", func(ctx SpecContext) {
			err := valjean.Checkout(ctx, library, "Les Miserables")
			Expect(err).To(MatchError("Les Miserables is not in the library catalog"))
		}, SpecTimeout(time.Second*5))
	})
})
