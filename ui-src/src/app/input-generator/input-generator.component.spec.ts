import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InputGeneratorComponent } from './input-generator.component';

describe('InputGeneratorComponent', () => {
  let component: InputGeneratorComponent;
  let fixture: ComponentFixture<InputGeneratorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InputGeneratorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InputGeneratorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
