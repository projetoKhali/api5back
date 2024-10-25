import pandas as pd
from sqlalchemy import create_engine
from transformation import *
import pg8000

def create_db_connection():
    """Cria uma conexão com o banco de dados PostgreSQL usando SQLAlchemy."""
    connection_string = "postgresql+pg8000://postgres:pass_warehouse@localhost:5433/postgres"
    return create_engine(connection_string)

def insert_dim_datetime(engine, df_dim_datetime):
    """Insere dados na dimensão datetime."""
    try:
        df_dim_datetime.to_sql(
            'dim_datetime',
            engine,
            if_exists='append',
            index=False,
            method='multi',
            chunksize=1000
        )
        print("Dados inseridos com sucesso em dim_datetime")
    except Exception as e:
        print(f"Erro ao inserir em dim_datetime: {str(e)}")
        raise

def insert_dim_user(engine, df_dim_user):
    """Insere dados na dimensão user."""
    try:
        # Renomeando as colunas para corresponder à estrutura do banco
        df_to_insert = df_dim_user.rename(columns={
            'id': 'dim_usr_id',
            'name': 'dim_usr_name',
            'occupation': 'dim_usr_ocupation'
        })
        
        df_to_insert.to_sql(
            'dim_user',
            engine,
            if_exists='append',
            index=False,
            method='multi',
            chunksize=1000
        )
        print("Dados inseridos com sucesso em dim_user")
    except Exception as e:
        print(f"Erro ao inserir em dim_user: {str(e)}")
        raise

def insert_dim_process(engine, df_dim_process):
    """Insere dados na dimensão process."""
    try:
        # Renomeando as colunas para corresponder à estrutura do banco
        df_to_insert = df_dim_process.rename(columns={
            'id': 'dim_pc_id',
            'title': 'dim_pc_title',
            'initial_date': 'dim_pc_initial_date',
            'finish_date': 'dim_pc_finish_date',
            'status': 'dim_pc_status',
            'dim_usr_id': 'dim_usr_id',
            'description': 'dim_pc_description'
        })
        
        df_to_insert.to_sql(
            'dim_process',
            engine,
            if_exists='append',
            index=False,
            method='multi',
            chunksize=1000
        )
        print("Dados inseridos com sucesso em dim_process")
    except Exception as e:
        print(f"Erro ao inserir em dim_process: {str(e)}")
        raise

def insert_dim_vacancy(engine, df_dim_vacancy):
    """Insere dados na dimensão vacancy."""
    try:
        # Garantindo apenas as colunas necessárias
        columns = [
            'vc_id', 'vc_title', 'vc_num_positions',
            'vc_status', 'vc_location', 'usr_id',
            'vc_opening_date', 'vc_closing_date'
        ]
        
        df_to_insert = df_dim_vacancy[columns]
        
        df_to_insert.to_sql(
            'dim_vacancy',
            engine,
            if_exists='append',
            index=False,
            method='multi',
            chunksize=1000
        )
        print("Dados inseridos com sucesso em dim_vacancy")
    except Exception as e:
        print(f"Erro ao inserir em dim_vacancy: {str(e)}")
        raise

def insert_fact_hiring_process(engine, fact_hiring_process):
    """Insere dados na tabela fato hiring_process."""
    try:
        # Renomeando as colunas para corresponder à estrutura do banco
        df_to_insert = fact_hiring_process.rename(columns={
            'id': 'fac_id',
            'dim_process_id': 'dim_process_dim_pc_id',
            'dim_vacancy_id': 'dim_vacancy_vc_id',
            'dim_user_id': 'dim_user_dim_usr_id',
            'dim_date_id': 'dim_datetime_dim_datetime_id',
            'met_sum_duration_hiring_process': 'met_sum_duration_hiring_proces'
        })
        
        columns = [
            'fac_id', 'dim_process_dim_pc_id', 'dim_vacancy_vc_id',
            'dim_user_dim_usr_id', 'dim_datetime_dim_datetime_id',
            'met_total_candidates_applied', 'met_total_candidates_interviewed',
            'met_total_candidates_hired', 'met_sum_duration_hiring_proces',
            'met_sum_salary_initial', 'met_total_feedback_positive',
            'met_total_neutral', 'met_total_negative'
        ]
        
        df_to_insert = df_to_insert[columns]
        
        df_to_insert.to_sql(
            'fact_hiring_process',
            engine,
            if_exists='append',
            index=False,
            method='multi',
            chunksize=1000
        )
        print("Dados inseridos com sucesso em fact_hiring_process")
    except Exception as e:
        print(f"Erro ao inserir em fact_hiring_process: {str(e)}")
        raise

def load_data_to_postgres(df_dim_datetime, df_dim_user, df_dim_process, 
                         df_dim_vacancy, fact_hiring_process):
    """Função principal que coordena a inserção de todos os dados."""
    try:
        engine = create_db_connection()
        
        print("Iniciando inserção dos dados...")
        
        print("\nInserindo dados na dimensão datetime...")
        insert_dim_datetime(engine, df_dim_datetime)
        
        print("\nInserindo dados na dimensão user...")
        insert_dim_user(engine, df_dim_user)
        
        print("\nInserindo dados na dimensão process...")
        insert_dim_process(engine, df_dim_process)
        
        print("\nInserindo dados na dimensão vacancy...")
        insert_dim_vacancy(engine, df_dim_vacancy)
        
        print("\nInserindo dados na tabela fato hiring_process...")
        insert_fact_hiring_process(engine, fact_hiring_process)
        
        print("\nTodos os dados foram inseridos com sucesso!")
        
    except Exception as e:
        print(f"\nErro durante a inserção dos dados: {str(e)}")
        raise e
    
    finally:
        if engine:
            engine.dispose()

if __name__ == "__main__":
    # Chamada da função principal com os DataFrames já criados
    load_data_to_postgres(
        dim_datetime,
        df_dim_user,
        df_dim_process,
        df_dim_vacancy,
        fact_hiring_process
    )