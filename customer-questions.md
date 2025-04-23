# Question 1

Hello,

I'm new to search engines, and there are a lot of concepts I'm not educated on. To make my onboarding smoother, it'd help if you could provide me with some definitions of the following concepts:

- Records
- Indexing

I'm also struggling with understanding what types of metrics would be useful to include in the "Custom Ranking."

Cheers,
George

---

Hi George,

Welcome aboard! Don’t worry, I’m here to help you get up to speed quickly.

Here are the definitions you asked for:

- **Records**: These are the individual items you want your users to find, like products, restaurants, or recipes. Each record is made up of *attributes* (key-value pairs) that describe it.
  For example, a recipe record might look like:
  ```json
    { "title": "Blackberry and blueberry pie", "likes": 1128, "gluten_free": false }
  ```
- **Indexing**: This is the process of organizing and storing your records so they’re quick and easy to search. Algolia creates an index from your data to speed up searches. 
  It is like creating a quick-reference guide for books in a library, helping you find a specific book fast. Instead of checking every shelf, it points you straight to the right spot.

**Custom Ranking** metrics help prioritize results based on what matters to your users. Useful examples include numerical values like star ratings, number of reviews, price, or distance, which ensure popular or relevant items rank higher in search results.
So for a restaurant locator app, if users type “steak”, a custom ranking with star ratings puts the most popular steakhouses at the top of the results.

Please feel free to reach out if you need more help.

Regards,
Marco

P.S. For more details, you can check out these links:  
- [Glossary: Records](https://www.algolia.com/doc/glossary/#record)  
- [Preparing Records](https://www.algolia.com/doc/guides/sending-and-managing-data/prepare-your-data/#algolia-records)  
- [Glossary: Indexing](https://www.algolia.com/doc/glossary/#index)  
- [Sending Data (incl. Indexing)](https://www.algolia.com/doc/guides/sending-and-managing-data/send-and-update-your-data/)  
- [Custom Ranking](https://www.algolia.com/doc/guides/managing-results/must-do/custom-ranking/)

# Question 2

Hello,

Sorry to give you the kind of feedback that I know you do not want to hear, but I really hate the new dashboard design. Clearing and deleting indexes are now several clicks away. I am needing to use these features while iterating, so this is inconvenient.

Regards,
Matt

---

Hi Matt,

I’m really sorry that you’re frustrated with the new dashboard design. I totally get how extra clicks for clearing and deleting indexes can be inconvenient, especially while iterating.

Thank you for sharing this; we appreciate your feedback. I’ll pass it along to the appropriate teams to look into streamlining those actions.

In the meantime, you might find it faster to use our API for these tasks. You can check out this guide for details: [API Reference](https://www.algolia.com/doc/api-reference/api-methods/delete-index/).

If there’s anything else I can do to help, please let me know.

Regards,
Marco

# Question 3

Hi,

I'm looking to integrate Algolia in my website. Will this be a lot of development work for me? What's the high level process look like?

Regards,
Leo

---

Hi Leo,

Happy to know that you’re considering Algolia for your website! It’s a great choice, and we’ve made it developer-friendly to keep things smooth.

For a basic integration, it’s pretty straightforward. In fact, it is possible to create an Algolia interface in just two minutes. 

Here’s the high-level process:

1. **Set Up an Algolia Account**: Sign up and get your API keys.
2. **Prepare Your Data**: Format, upload, and configure your data via API or dashboard.  
3. **Install the Library**: Add Algolia’s search library (e.g., via JavaScript).  
4. **Configure Search**: Customize the UI and styling.  
5. **Test and Deploy**: Test it out, then go live!

We’ve got plenty of resources like docs and tutorials to help. You can try our [two-minute onboarding demo](https://www.algolia.com/doc/onboarding/) or the [Quick Start guide](https://www.algolia.com/doc/guides/getting-started/quick-start/). 

If you would like tailored advice, please feel free to share more about your project.

Regards,
Marco
