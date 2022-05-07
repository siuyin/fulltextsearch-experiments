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

# Components
[![](https://mermaid.ink/img/pako:eNqFjkFrAjEQhf9KmNMWXHvfQ6GgB7H00C09DZQxmdVANlmSSbui_ndHsefeHu99M--dwCbH0MEQ0q89UBbz9oHRJduskkWMkWUIfv4WL4HL0pafJ4zezc0mOp4V2IW6Z_WKNL2kzGq9v372ZsvH9otCvWX3C9O2L-dSpyl4zuZZOVuLpJHz2Wjf7cM_CCxA5Uje6eATRmMQ5MAjI3QqHQ9UgyBgvChaJ0fCa-d1FXQDhcILoCqpP0YLneTKf9DK0z7T-KAuV5qrY7k)](https://mermaid.live/edit#pako:eNqFjkFrAjEQhf9KmNMWXHvfQ6GgB7H00C09DZQxmdVANlmSSbui_ndHsefeHu99M--dwCbH0MEQ0q89UBbz9oHRJduskkWMkWUIfv4WL4HL0pafJ4zezc0mOp4V2IW6Z_WKNL2kzGq9v372ZsvH9otCvWX3C9O2L-dSpyl4zuZZOVuLpJHz2Wjf7cM_CCxA5Uje6eATRmMQ5MAjI3QqHQ9UgyBgvChaJ0fCa-d1FXQDhcILoCqpP0YLneTKf9DK0z7T-KAuV5qrY7k)

# Executables
`grep -ir dflt.Env` to see all environment variable configurables.

## Index building
`cmd/bluge-idx/main.go`: creates the index

## Searching
`cmd/bluge-search/main.go`: searches the index

Eg. bluge-search 'Title:star trek Description:picard'  
The above searches for "star trek" in the title or "picard" in the description.

Eg. bluge-search 'Title:star trek +Description:picard'  
The above searches for "star trek" in the title AND "picard" in the description.

For additional syntax details, see: https://blevesearch.com/docs/Query-String-Query/
