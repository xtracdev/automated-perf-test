import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ApplicationPropertiesComponent } from './application-properties.component';

describe('ApplicationPropertiesComponent', () => {
  let component: ApplicationPropertiesComponent;
  let fixture: ComponentFixture<ApplicationPropertiesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ApplicationPropertiesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ApplicationPropertiesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
