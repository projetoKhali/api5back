import pandas as pd
from sqlalchemy import create_engine, text, Column, Integer, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import logging
from datetime import datetime
from transformation import *

class PostgreSQLLoader:
    def __init__(self):
        """Inicializa o loader com as credenciais do PostgreSQL."""
        self.connection_string = 'postgresql+pg8000://postgres:pass_warehouse@localhost:5433/postgres'
        self.engine = create_engine(self.connection_string)
        self.setup_logging()

    def setup_logging(self):
        """Configura o logging para monitorar o processo de carga."""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler(f'etl_loader_{datetime.now().strftime("%Y%m%d_%H%M%S")}.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)

    def test_connection(self):
        """Testa a conexão com o banco de dados."""
        try:
            with self.engine.connect() as conn:
                self.logger.info("Conexão com o banco de dados estabelecida com sucesso!")
                return True
        except Exception as e:
            self.logger.error(f"Erro ao conectar com o banco de dados: {str(e)}")
            return False

    def get_dim_datetime_id(self):
        """Obtém o id do registro de dim_datetime mais recente."""
        query = """
            SELECT id
            FROM dim_datetime
            ORDER BY id DESC
            LIMIT 1;
        """
        with self.engine.connect() as conn:
            result = conn.execute(text(query))
            dim_datetime_id = result.scalar()
            self.logger.info(f"Obteve id de dim_datetime: {dim_datetime_id}")
            return dim_datetime_id

    def load_table(self, df, table_name):
        """Carrega dados em uma tabela."""
        try:
            with self.engine.begin() as conn:
                df.to_sql(
                    name=table_name,
                    con=conn,
                    if_exists='append',
                    index=False,
                    method='multi',
                    chunksize=1000
                )
                self.logger.info(f"Carregados {len(df)} registros na tabela {table_name}")
        except Exception as e:
            self.logger.error(f"Erro ao carregar tabela {table_name}: {str(e)}")
            raise

    def load_fact_table(self, fact_table):
        """Carrega dados na tabela fato."""
        dim_datetime_id = self.get_dim_datetime_id()
        fact_table['dim_date_id'] = dim_datetime_id
        self.load_table(fact_table, 'fact_hiring_process')

    def main(self):
        """Função principal para teste de inserção."""
        try:
            loader = PostgreSQLLoader()
            if loader.test_connection():
                print("Conexão bem-sucedida! Iniciando carga de teste...")
                loader.load_table(df_dim_user, 'dim_user')
                loader.load_table(df_dim_process, 'dim_process')
                loader.load_table(df_dim_datetime, 'dim_datetime')
                loader.load_table(df_dim_vacancy, 'dim_vacancy')
                loader.load_fact_table(fact_hiring_process)
                print("Carga de teste concluída com sucesso!")
            else:
                print("Não foi possível estabelecer conexão com o banco de dados.")
        except Exception as e:
            logging.error(f"Erro no processo de teste: {str(e)}")
            raise

if __name__ == "__main__":
    PostgreSQLLoader().main()