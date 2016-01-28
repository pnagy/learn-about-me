package main

import (
    "github.com/jinzhu/gorm"
)

type User struct {
    gorm.Model
    Username       string `sql:"unique_index"`
    Name           string
    Introduction   string `sql:"type:varchar(1000);unique"`
    ProfilePicture string
    Email          string
    BirthDate      string
    Location       string
    Contacts       []Contact
    Schools        []School
    Skills         []Skill
}

type Contact struct {
    gorm.Model
    UserID uint `sql:"index"`
    Key    string
    Value  string
    Weight uint
}

type School struct {
    gorm.Model
    UserID uint `sql:"index"`
    Name   string
    City   string
    From   uint
    To     uint
}

type Job struct {
    gorm.Model
    UserID      uint `sql:"index"`
    Company     string
    Position    string
    From        uint
    To          uint
    Description string
}

type Skill struct {
    gorm.Model
    UserID uint `sql:"index"`
    Name   string
    Level  uint
}
