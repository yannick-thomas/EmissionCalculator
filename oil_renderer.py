import tkinter as tk
from tkinter import ttk
import customtkinter

from renderer import Renderer
from oil_calc import OilCalculation

class OilRenderer(Renderer):
    is_valid_calc = True

    def build_base(self):
        self.quantity = tk.StringVar(self.window)
        self.main_content = customtkinter.CTkFrame(self.window, corner_radius=0, fg_color="#232323")
        self.main_content.grid(row=0, column=1, rowspan=5, sticky="nsew")
        self.main_content.grid_rowconfigure(5, weight=1)
        main_content_heading = customtkinter.CTkLabel(
            self.main_content, 
            text="Emissionsberechnung für Heizöl", 
            font=customtkinter.CTkFont(size=20, weight="bold")
        )
        main_content_heading.grid(row=1, column=0, padx=(20, 0), pady=(20, 0), columnspan=4, sticky="nw")
        label_quantity_oil= customtkinter.CTkLabel(self.main_content, text="Liefermenge:", font=customtkinter.CTkFont(size=16))
        label_quantity_oil.grid(row=2, column=0,columnspan=1, padx=(20, 0), pady=(20, 0), sticky="nw")
        entry_quantity_oil = customtkinter.CTkEntry(
            self.main_content, 
            width=50, 
            font=customtkinter.CTkFont(size=12, weight="bold"), 
            textvariable=self.quantity
        )
        entry_quantity_oil.grid(row=2, column=1, padx=(155, 0), pady=(20, 0), sticky="nw")
        label_quantity_measure = customtkinter.CTkLabel(
            self.main_content, text="l", 
            font=customtkinter.CTkFont(size=16, weight="bold")
        )
        label_quantity_measure.grid(row=2, column=2, padx=(2, 0), pady=(20, 0), sticky="nw")
        calc_button_oil = customtkinter.CTkButton(
            self.main_content, 
            text="Berechnen", 
            width=120, 
            command=self.calc
        )
        calc_button_oil.grid(row=2, column=5, padx=(20, 0), pady=(20, 0))
        emissions_heading_label = customtkinter.CTkLabel(
            self.main_content,
            font=customtkinter.CTkFont(size=16),
            text_color="grey"
        )
        emissions_heading_label.grid(row=3,columnspan=3, column=0 ,padx=(20, 0), pady=(30, 0), sticky="nw")
        emissions_formula = customtkinter.CTkLabel(
            self.main_content,
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        emissions_formula.grid(row=4 ,columnspan=6,column=0, padx=(20, 0), pady=(0, 0), sticky="nw")
        self.label_calc_failure = customtkinter.CTkLabel(
            self.main_content, 
            font=customtkinter.CTkFont(size=16, weight="bold")
        )
        self.label_calc_failure.grid(row=3, column=1, padx=(100, 0), pady=(30, 0), sticky="nw")
        self.window.bind('<Return>', lambda event: self.calc())

        
    def calc(self):
        calculation = OilCalculation(self.quantity.get())
        result = calculation.calculate()
        self.render_post_calculation(result)

    def render_post_calculation(self, result):

        emissions_heading_label.grid_forget()
        emissions_formula.grid_forget()
        emissions_equals.grid_forget()
        emissions_result.grid_forget()
        emission_component_heading.grid_forget()
        emission_component_formula.grid_forget()
        emission_component_equals.grid_forget()
        emission_component_result_label.grid_forget()
        self.label_calc_failure.grid_forget()

        emissions_heading_label = customtkinter.CTkLabel(
            self.main_content,
            text="Brennstroffemissionen der Heizöllieferung:",
            font=customtkinter.CTkFont(size=16),
            text_color="grey"
        )
        emissions_heading_label.grid(row=3,columnspan=3, column=0 ,padx=(20, 0), pady=(30, 0), sticky="nw")
        emissions_formula = customtkinter.CTkLabel(
            self.main_content,
            text="42,8 GJ/t x 0,074 t CO2/GJ x 0,845t/1000l x " + self.quantity.get()+ "l", 
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        emissions_formula.grid(row=4 ,columnspan=6,column=0, padx=(20, 0), pady=(0, 0), sticky="nw")
        emissions_equals = customtkinter.CTkLabel(self.main_content, text="=")
        emissions_equals.grid(row=4,columnspan=6, column=2, padx=(50, 0), pady=(0, 0), sticky="nw")
        emissions_result = customtkinter.CTkLabel(
            self.main_content,
            text=str(round(result[0], 2)).replace(".", ",") + "t CO2",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        emissions_result.grid(row=4,columnspan=6, column=3, padx=(10, 0), pady=(0, 0), sticky="ne")
        emission_component_heading = customtkinter.CTkLabel(
            self.main_content, 
            text="Preisbestandteil CO2 Kosten (inkl. Umsatzsteuer):",
            font=customtkinter.CTkFont(size=16),
            text_color="grey"
        )
        emission_component_heading.grid(row=5, columnspan= 6, column=0, padx=(20, 0), pady=(20, 0), sticky="nw")
        emission_component_formula = customtkinter.CTkLabel(
            self.main_content,
            text= str(round(result[0], 2)).replace(".", ",") + "t CO2 x 30€/t CO2 x 1,12 MwSt.",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        emission_component_formula.grid(row=6, columnspan=6, column=0, padx=(20, 0), pady=(0, 460), sticky="nw")
        emission_component_equals = customtkinter.CTkLabel(
            self.main_content,
            text="="
        )
        emission_component_equals.grid(row=6,columnspan=6, column=2, padx=(50, 0), pady=(0, 0), sticky="nw")
        # emissions = round(RESULT[0], 2)
        emission_component_result_label = customtkinter.CTkLabel(
            self.main_content,
            text=str(round(result[1], 2)) + "€",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        emission_component_result_label.grid(row=6, columnspan=6, column=3, padx=(0, 20), pady=(0, 0), sticky="ne")
