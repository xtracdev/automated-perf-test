import {async, fakeAsync, ComponentFixture, TestBed} from "@angular/core/testing";
import {ConfigurationsComponent} from "./configurations.component";

import {CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA} from "@angular/core";
import {ToastModule} from "ng2-toastr/ng2-toastr";
import {ConfigurationService} from "./configuration.service";
import {TestSuiteService} from "../test-suites/test-suite.service";
import {HttpClientModule} from "@angular/common/http";
import {FormsModule} from "@angular/forms";
import {ToastsManager, ToastOptions} from "ng2-toastr/ng2-toastr";
import {Observable} from "rxjs/Observable";
import "rxjs/add/observable/of";
import "rxjs/add/observable/throw";

describe("ConfigurationsComponent", () => {
  let component: ConfigurationsComponent;
  let fixture: ComponentFixture<ConfigurationsComponent>;
  let configurationService, toastr;

  beforeEach(
    async(() => {
      TestBed.configureTestingModule({
        providers: [ConfigurationService, ToastsManager, ToastOptions, TestSuiteService],
        declarations: [ConfigurationsComponent],
        imports: [FormsModule, HttpClientModule, ToastModule.forRoot()],
        schemas: [NO_ERRORS_SCHEMA,
          CUSTOM_ELEMENTS_SCHEMA]
      }).compileComponents();

      fixture = TestBed.createComponent(ConfigurationsComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();

      configurationService = TestBed.get(ConfigurationService);
      toastr = TestBed.get(ToastsManager);
    }));

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  describe("post config file", () => {
    it("should submit config data if it is valid", () => {
      spyOn(configurationService, "postConfig$").and.returnValue(Observable.of("some value"));
      spyOn(toastr, "success");
      component.configPath = "";
      component.onSubmit({});
      fixture.detectChanges();
      expect(toastr.success).toHaveBeenCalledWith("Your Data has Been Saved!", "Success!");
    });

    it("should reject and show a error toast should display", () => {
      spyOn(configurationService, "postConfig$").and.returnValue(Observable.throw("some value"));
      spyOn(toastr, "error");
      component.configPath = "";
      component.onSubmit({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("Your Data did not Save!", "An Error Occurred!");
    });

    it("should show a error toast when file not found", () => {
      spyOn(configurationService, "postConfig$").and.returnValue(Observable.throw({status: 409}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.onSubmit({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("File Already Exists!", "An Error Occurred!");
    });

    it("should show a error toast when there is a bad request", () => {
      spyOn(configurationService, "postConfig$").and.returnValue(Observable.throw({status: 400}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.onSubmit({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("Some of the Fields do not Conform to the Schema!", "An Error Occurred!");
    });

    it("should show a error toast when the serivce is down", () => {
      spyOn(configurationService, "postConfig$").and.returnValue(Observable.throw({status: 500}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.onSubmit({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("An Error has Occurred!", "Check the logs!");
    });
  });

  describe("get config file", () => {
    it("should get file config data if it is valid config path and filename", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.of("some value"));
      spyOn(toastr, "success");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.onGetFile();
      fixture.detectChanges();
      expect(toastr.success).toHaveBeenCalledWith("Success!");
      expect(component.formData).toEqual("some value");
    });

    it("should fail to get file", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.throw("some value"));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.onGetFile();
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("Your Data was Not Retrieved!", "An Error Occurred!");
    });

    it("should fail to get file when file not found", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.throw({status: 404}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.onGetFile();
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("File Not Found!", "An Error Occured!");
    });

    it("should fail to get file when there is a bad request", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.throw({status: 400}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.onGetFile();
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("Check Your Field Inputs", "An Error Occurred!");
    });

    it("should fail to get file when the serivce is down", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.throw({status: 500}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.onGetFile();
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("An Error has Occurred!", "Check the logs!");
    });
  });


  describe("update config file", () => {
    it("should update file config data if it is valid config path and filename", () => {
      spyOn(configurationService, "putConfig$").and.returnValue(Observable.of("some value"));
      spyOn(toastr, "success");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      this.formData = "";
      component.onUpdate({});
      fixture.detectChanges();
      expect(toastr.success).toHaveBeenCalledWith("Success!");
    });

    it("should fail to update file", () => {
      spyOn(configurationService, "putConfig$").and.returnValue(Observable.throw("some value"));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      this.formData = "";
      component.onUpdate({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("File Was Not Updated!", "An Error Occurred!");
    });

    it("should fail to update file when file not found", () => {
      spyOn(configurationService, "putConfig$").and.returnValue(Observable.throw({status: 404}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      this.formData = "";
      component.onUpdate({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("File Not Found", "An Error Occured!");
    });

    it("should fail to update file when there is a bad request", () => {
      spyOn(configurationService, "putConfig$").and.returnValue(Observable.throw({status: 400}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      this.formData = "";
      component.onUpdate({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("File Must be Specified!", "An Error Occurred!");
    });

    it("should fail to update file when the serivce is down", () => {
      spyOn(configurationService, "putConfig$").and.returnValue(Observable.throw({status: 500}));
      spyOn(toastr, "error");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      this.formData = "";
      component.onUpdate({});
      fixture.detectChanges();
      expect(toastr.error).toHaveBeenCalledWith("Internal Server Error!");
    });
  });


  describe("cancel form", () => {
    it("should get the previous loaded data", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.of("some value"));
      spyOn(toastr, "success");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.formData = "old data";

      component.onCancel();
      fixture.detectChanges();
      expect(toastr.success).toHaveBeenCalledWith("Previous Data Reloaded!");
      expect(component.formData).toEqual("some value");
    });

    it("should clear alll data when error occurs", () => {
      spyOn(configurationService, "getConfig$").and.returnValue(Observable.throw("some value"));
      spyOn(toastr, "success");
      component.configPath = "";
      component.xmlFileName = "fileName.xml";
      component.formData = "old data";

      component.onCancel();
      fixture.detectChanges();
      expect(component.configPath).toBeUndefined();
      expect(component.xmlFileName).toBeUndefined();
      expect(component.formData).toBeUndefined();
    });
  });
});
