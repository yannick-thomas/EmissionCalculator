import tkinter as tk
from tkinter import ttk
import customtkinter

from calculator.renderer.renderer import Renderer
from calculator.calc.briquettes_calc import BriquettesCalculation
from calculator.misc.print_handler import PrintHandler

class BriquettesRenderer(Renderer):
    is_valid_calc = ''

    def build_base(self):
        self.quantity = tk.StringVar(self.window)
        self.main_content = customtkinter.CTkFrame(self.window, corner_radius=0, fg_color="#242424")
        # self.main_content = customtkinter.CTkFrame(self.window, corner_radius=0)
        self.main_content.grid(row=0, column=1, rowspan=5, sticky="nsew")
        self.main_content.grid_rowconfigure(5, weight=1)
        main_content_heading = customtkinter.CTkLabel(
            self.main_content, 
            text="Emissionsberechnung für Briketts", 
            font=customtkinter.CTkFont(size=20, weight="bold"),
            text_color=self.font_color
        )
        main_content_heading.grid(row=1, column=0, padx=(20, 0), pady=(20, 0), columnspan=4, sticky="nw")
        label_quantity_oil= customtkinter.CTkLabel(self.main_content, text="Liefermenge:", font=customtkinter.CTkFont(size=16))
        label_quantity_oil.grid(row=2, column=0,columnspan=1, padx=(20, 0), pady=(20, 0), sticky="nw")
        entry_quantity_oil = customtkinter.CTkEntry(
            self.main_content, 
            width=70, 
            font=customtkinter.CTkFont(size=12, weight="bold"), 
            textvariable=self.quantity,
            text_color=self.font_color
        )
        entry_quantity_oil.grid(row=2, column=1, padx=(155, 0), pady=(20, 0), sticky="nw")
        label_quantity_measure = customtkinter.CTkLabel(
            self.main_content, text="t", 
            font=customtkinter.CTkFont(size=16, weight="bold"),
            text_color=self.font_color
        )
        label_quantity_measure.grid(row=2, column=2, padx=(2, 0), pady=(20, 0), sticky="nw")
        calc_button_oil = customtkinter.CTkButton(
            self.main_content, 
            text="Berechnen", 
            width=120, 
            command=self.calc,
            text_color=self.font_color,
            fg_color=self.button_color,
            hover_color=self.hover_color
        )
        calc_button_oil.grid(row=2, column=5, padx=(80, 0), pady=(20, 0))
        self.label_calc_failure = customtkinter.CTkLabel(
            self.main_content,
            text = "Fehler!",
            font=customtkinter.CTkFont(size=16, weight="bold"),
            text_color=self.font_color
        )
        self.window.bind('<Return>', lambda event: self.calc())

        
    def calc(self):
        calculation = BriquettesCalculation(self.quantity.get())
        self.result = calculation.calculate()
        self.render_post_calculation(self.result)

    def render_post_calculation(self, result):

        def print_label():
            handler = PrintHandler(self.result)
            handler.print_result()

        if not result[0]:
            if not self.is_valid_calc == '':
                self.fuel_emission_calc_label.grid_forget()
                self.fuel_emission_calc_formula.grid_forget()
                self.fuel_emission_calc_equals.grid_forget()
                self.fuel_emission_calc_equals_result.grid_forget()
                self.price_comp_co2_costs_formula.grid_forget()
                self.price_comp_co2_costs_equals.grid_forget()
                self.price_comp_co2_costs_label.grid_forget()
                self.result_co2_costs_label.grid_forget()
                self.energy_content_formula.grid_forget()
                self.energy_content_equals.grid_forget()
                self.energy_content_label.grid_forget()
                self.result_energy_content_label.grid_forget()
                self.button_print.grid_forget()
                self.label_calc_failure.grid_forget()
            self.label_calc_failure.grid(row=3, column=1, padx=(100, 0), pady=(30, 0), sticky="nw")
            self.is_valid_calc = False
            return

        if self.is_valid_calc != True:
            self.fuel_emission_calc_label = customtkinter.CTkLabel(
                self.main_content,
                text="Brennstoffemissionen der aktuellen Lieferung:", 
                font=customtkinter.CTkFont(size=16), 
                text_color=self.header_color
            )
            self.fuel_emission_calc_formula = customtkinter.CTkLabel(
                self.main_content,
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.formula_color
            )
            self.fuel_emission_calc_equals = customtkinter.CTkLabel(
                self.main_content,
                text="≈",
                font=customtkinter.CTkFont(size=18),
                text_color=self.formula_color
            )
            self.fuel_emission_calc_equals_result = customtkinter.CTkLabel(
                self.main_content,
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.formula_color
            )
            self.price_comp_co2_costs_label = customtkinter.CTkLabel(
                self.main_content, 
                text="Preisbestandteil CO2 Kosten:", 
                font=customtkinter.CTkFont(size=16), 
                text_color=self.header_color
            )
            self.price_comp_co2_costs_formula = customtkinter.CTkLabel(
                self.main_content, 
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.formula_color
            )
            self.price_comp_co2_costs_equals = customtkinter.CTkLabel(
                self.main_content, 
                text="=", 
                font=customtkinter.CTkFont(size=18),
                text_color=self.formula_color
            )
            self.result_co2_costs_label = customtkinter.CTkLabel(
                self.main_content,
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.formula_color
            )
            self.energy_content_label = customtkinter.CTkLabel(
                self.main_content, text="Energiegehalt der Lieferung:",
                font=customtkinter.CTkFont(size=16),
                text_color=self.header_color
            )
            self.energy_content_formula = customtkinter.CTkLabel(
                self.main_content,
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.formula_color
            )
            self.energy_content_equals = customtkinter.CTkLabel(
                self.main_content,
                text="=",
                font=customtkinter.CTkFont(size=18),
                text_color=self.formula_color
            )
            self.result_energy_content_label = customtkinter.CTkLabel(
                self.main_content,
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.formula_color
            )
            self.button_print = customtkinter.CTkButton(
                self.main_content,
                text="Drucken",
                command=print_label,
                font=customtkinter.CTkFont(size=14, weight="bold"),
                text_color=self.font_color,
                fg_color=self.button_color,
                hover_color=self.hover_color
            )

        self.fuel_emission_calc_formula.configure(
            text="19,0 GJ/t x 0,0992t CO2 / GJ x 1 x " + self.quantity.get() + "t", 
        )
        self.fuel_emission_calc_equals_result.configure(
            text=result[1] + "kg CO2"
        )
        self.price_comp_co2_costs_formula.configure(
            text=result[1] + "kg CO2 x 30€/ CO2 t x 1,19", 
        )
        self.result_co2_costs_label.configure(
            text=result[2] + " €",
        )
        self.energy_content_formula.configure(
            text="19,0 GJ x 1 x " + self.quantity.get() + "t"
        )
        self.result_energy_content_label.configure(
            text=result[3] + " kWh"
        )

        self.energy_content_equals.grid(row=6,columnspan=3, column=0, padx=(330, 0), pady=(65, 0), sticky="nw")
        
        self.fuel_emission_calc_label.grid(row=3,columnspan=3, column=0 ,padx=(20, 0), pady=(30, 0), sticky="nw")
        self.fuel_emission_calc_formula.grid(row=4,columnspan=3, column=0, padx=(20, 0), pady=(0, 0), sticky="nw")
        self.fuel_emission_calc_equals.grid(row=4,columnspan=3, column=0, padx=(330, 0), pady=(0, 0), sticky="nw")
        self.fuel_emission_calc_equals_result.grid(row=4,columnspan=3, column=3, padx=(20, 0), pady=(0, 0), sticky="nw")
        
        self.price_comp_co2_costs_label.grid(row=5, columnspan=3 ,column=0,padx=(20, 0), pady=(10, 0), sticky="nw")
        self.price_comp_co2_costs_formula.grid(row=6, columnspan=3 ,column=0,padx=(20, 0), pady=(0, 270), sticky="nw")
        self.price_comp_co2_costs_equals.grid(row=6,columnspan=3, column=0, padx=(330, 0), pady=(0, 0), sticky="nw") 
        self.result_co2_costs_label.grid(row=6,columnspan=3, column=3, padx=(20, 0), pady=(0, 0), sticky="nw")
        
        self.energy_content_label.grid(row=6,columnspan=3, column=0, padx=(20, 0), pady=(40, 0), sticky="nw")
        self.energy_content_formula.grid(row=6,columnspan=3, column=0, padx=(20, 0), pady=(65, 0), sticky="nw")
        self.result_energy_content_label.grid(row=6,columnspan=3, column=3, padx=(20, 0), pady=(65, 0), sticky="nw")
        self.button_print.grid(row=6,columnspan=3, column=1, padx=(50, 0), pady=(120, 0), sticky="nw")

        self.is_valid_calc = True
