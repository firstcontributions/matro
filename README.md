
<p align="center">
  
  <img src="./assets/logo.png" width="200" />
  <br/>
  <p align="center">
    <b>Matro </b> is a graphql code generator inspired by swagger
  </p>
  <!-- <p align="center">
    A <a href="">Firstcontributions</b> initiative
  </p> -->
  <p align="center">
    <a href= "https://github.com/firstcontributions/matro/actions/"> 
        <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/firstcontributions/matro/CI">
    </a>
    <a href= "https://github.com/firstcontributions/matro/issues"> 
        <img alt="GitHub issues" src="https://img.shields.io/github/issues/firstcontributions/matro">
    </a>
    <a href= "https://github.com/firstcontributions/matro/blob/main/LICENSE"> 
        <img alt="GitHub license" src="https://img.shields.io/github/license/firstcontributions/matro">
    </a>
    <a href= "https://codeclimate.com/github/firstcontributions/matro/maintainability"> 
        <img alt="Code Climate maintainability" src="https://api.codeclimate.com/v1/badges/99dfc661e165766b7528/maintainability">
    </a>
    <a href= ""> 
        <img alt="Code Climate technical debt" src="https://img.shields.io/codeclimate/tech-debt/firstcontributions/matro">
    </a>
    <a href= "https://codeclimate.com/github/firstcontributions/matro/test_coverage"> 
        <img alt="Code Climate coverage" src="https://api.codeclimate.com/v1/badges/99dfc661e165766b7528/test_coverage">
    </a>
    <a href="https://deepsource.io/gh/firstcontributions/matro/?ref=repository-badge}" target="_blank">
      <img alt="DeepSource" title="DeepSource" src="https://deepsource.io/gh/firstcontributions/matro.svg/?label=active+issues&show_trend=truetoken=CyxagxqXgW4t86z6c-IDpfy7"/>
    </a>
    <a href="https://deepsource.io/gh/firstcontributions/matro/?ref=repository-badge}" target="_blank">
      <img alt="DeepSource" title="DeepSource" src="https://deepsource.io/gh/firstcontributions/matro.svg/?label=resolved+issues&show_trend=true&token=CyxagxqXgW4t86z6c-IDpfy7"/>
    </a>
  </p>
  
</p>


## Getting started

### 1. Pre-requisites
  1. Go version 17.x+
  2. For grpc
      1. google.golang.org/grpc v1.44.0 
      2. google.golang.org/protobuf v1.27.1
  3. Gnu Make

### 2. How to Install
```sh
  $. git clone git@github.com:firstcontributions/matro.git
  $. cd matro
  $. make config
  $. make build
```

## Progress
---
- [x] GraphQL Schema generator
- [x] gRPC proto buf & service stub generator
- [x] gRPC service implementation
- [x] Models/ Store generator
- [x] Graphql golang implementation
- [ ] React-relay generator 
