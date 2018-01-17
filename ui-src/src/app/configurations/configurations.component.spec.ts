import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import { Component, OnInit } from "@angular/core";
import { Http } from "@angular/http";
import { AutomatedUIServices } from "../automated-ui-services";
import { ConfigurationsComponent } from "./configurations.component";
import { JsonSchemaFormModule } from "angular2-json-schema-form";
import { HttpClientModule, HttpClient } from "@angular/common/http";

describe("ConfigurationsComponent", () => {
  let component: ConfigurationsComponent;
  let fixture: ComponentFixture<ConfigurationsComponent>;

  beforeEach(
    async(() => {
      TestBed.configureTestingModule({
        providers: [AutomatedUIServices],
        declarations: [ConfigurationsComponent],
        imports: [JsonSchemaFormModule, HttpClientModule]
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigurationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
