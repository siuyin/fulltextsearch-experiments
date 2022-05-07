# Full text search experiments
The dataset: netflix_titles.csv is from kaggle: https://www.kaggle.com/datasets/yogithasatyasri/netflix-shows-exploratory-analysis .

# Requirements

## Feature: Full text search on a database
As a person looking for shows to watch on Netflix,
I want to be able to search Netflix's database of shows
So that I may discovery shows I was previously unaware of.

### Scenario: Search all public facing fields for a given search term
Given I am a Search user  
And I have a database of shows  
When I provide a search term  
Then I should see a list of entries,
ordered most relevant first,
with my search term highlighted in each entry.

# Executables
`cmd/bluge-idx/main.go`: creates the index

`cmd/bluge-search/main.go`: searches the index
